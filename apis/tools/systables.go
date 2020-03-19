package tools

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/models/tools"
	"go-admin/pkg"
	"go-admin/utils"
	"net/http"
	"strings"
)

// @Summary 分页列表数据
// @Description 生成表分页列表
// @Tags 工具 - 生成表
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/page [get]
func GetSysTableList(c *gin.Context) {
	var data tools.SysTables
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.TableName = c.Request.FormValue("tableName")
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.AssertErr(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageIndex"] = pageSize

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 工具 - 生成表
// @Param configKey path int true "configKey"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/info/{tableId} [get]
// @Security
func GetSysTables(c *gin.Context) {
	var data tools.SysTables
	data.TableId, _ = utils.StringToInt64(c.Param("tableId"))
	result, err := data.Get()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)

	var res models.Response
	res.Data = result
	mp := make(map[string]interface{})
	mp["rows"] = result.Columns
	mp["info"] = result
	res.Data = mp
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 添加表结构
// @Description 添加表结构
// @Tags 工具 - 生成表
// @Accept  application/json
// @Product application/json
// @Param tables query string false "tableName / 数据表名称"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sys/tables/info [post]
// @Security Bearer
func InsertSysTable(c *gin.Context) {
	var data tools.SysTables
	var dbTable tools.DBTables
	var dbColumn tools.DBColumns
	data.TableName = c.Request.FormValue("tables")
	data.CreateBy = utils.GetUserIdStr(c)

	dbTable.TableName = data.TableName
	dbtable, err := dbTable.Get()

	dbColumn.TableName = data.TableName
	dbcolumn, err := dbColumn.GetList()
	data.CreateTime = utils.GetCurrntTime()
	data.CreateBy = utils.GetUserIdStr(c)
	data.TableComment = dbtable.TableComment
	data.FunctionAuthor = "wenjianzhang"
	for i := 0; i < len(dbcolumn); i++ {
		var column tools.SysColumns
		column.ColumnComment = dbcolumn[i].ColumnComment
		column.ColumnName = dbcolumn[i].ColumnName
		column.ColumnType = dbcolumn[i].ColumnType
		column.Sort = utils.IntToString(i + 1)
		column.Insert = true
		column.IsInsert = "1"
		column.QueryType = "EQ"
		column.IsPk = "0"
		if strings.Contains(dbcolumn[i].ColumnKey, "PR") {
			column.IsPk = "1"
			column.Pk = true
		}
		column.IsRequired = "0"
		if strings.Contains(dbcolumn[i].IsNullable, "NO") {
			column.IsRequired = "1"
			column.Required = true
		}
		if strings.Contains(dbcolumn[i].ColumnType, "int") {
			column.GoType = "int64"
			column.HtmlType = "input"
		//} else if strings.Contains(dbcolumn[i].ColumnType, "char") {
		//	column.GoType = "bool"
		//	column.HtmlType = "input"
		} else {
			column.GoType = "string"
			column.HtmlType = "input"
		}

		data.Columns = append(data.Columns, column)
	}

	result, err := data.Create()
	pkg.AssertErr(err, "", -1)

	var res models.Response
	res.Data = result
	res.Msg = "添加成功！"
	c.JSON(http.StatusOK, res.ReturnOK())

}

// @Summary 删除表结构
// @Description 删除表结构
// @Tags 工具 - 生成表
// @Param tableId path int true "tableId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sys/tables/info/{tableId} [delete]
func DeleteSysTables(c *gin.Context) {
	var data tools.SysTables
	id, err := utils.StringToInt64(c.Param("tableId"))
	data.TableId = id
	_, err = data.Delete()
	pkg.AssertErr(err, "删除失败", 500)
	var res models.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}
