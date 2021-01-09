package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Customer struct {
		Name         string    `json:"cName"`
		Tel          uint64    `json:"cTel"`
		Address      string    `json:"cAddress"`
		ID           uint64    `json:"cID"`
		RegisterDate time.Time `json:"cRegisterDate"`
	}

	CustomerResponse struct {
		Name         string `json:"cName"`
		Tel          uint64 `json:"cTel"`
		Address      string `json:"cAddress"`
		ID           uint64 `json:"cID"`
		RegisterDate string `json:"cRegisterDate"`
		Message      string `json:"msg"`
	}

	AllCustomers struct {
		Size      uint64              `json:"size"`
		Customers []*CustomerResponse `json:"customers"`
		Message   string              `json:"msg"`
	}

	JustMSG struct {
		Message string `json:"msg"`
	}

	ReportResponse struct {
		TotalCustomers uint64 `json:"totalCustomers"`
		Period         int    `json:"period"`
		Message        string `json:"msg"`
	}
)

var (
	db          = map[uint64]*Customer{}
	last uint64 = 0
)

func getCustomersSlice(db map[uint64]*Customer) []*CustomerResponse {
	ans := make([]*CustomerResponse, 0, len(db))
	for _, c := range db {
		ans = append(ans, newCustomerResponse(c))
	}
	return ans
}

func newCustomerList(db map[uint64]*Customer) AllCustomers {
	return AllCustomers{
		Size:      uint64(len(db)),
		Customers: getCustomersSlice(db),
		Message:   "success",
	}
}

func newReportResponse(n uint64) ReportResponse {
	return ReportResponse{
		TotalCustomers: n,
		Period:         1,
		Message:        "success",
	}
}

func createError() JustMSG {
	return JustMSG{Message: "error"}
}

func timeToString(t time.Time) string {
	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
}

func newCustomerResponse(c *Customer) *CustomerResponse {
	return &CustomerResponse{
		Name:         c.Name,
		Tel:          c.Tel,
		Address:      c.Address,
		ID:           c.ID,
		RegisterDate: timeToString(c.RegisterDate),
		Message:      "success",
	}
}

func CreateNewCustomer(c echo.Context) error {

	newUser := &Customer{
		ID:           last + 1,
		RegisterDate: time.Now(),
	}
	if err := c.Bind(newUser); err != nil {
		return err
	}
	db[newUser.ID] = newUser
	last++

	response := newCustomerResponse(newUser)

	return c.JSON(http.StatusCreated, response)
}

func EditCustomer(c echo.Context) error {
	cidI, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		return err
	}
	cid := uint64(cidI)

	oldUser, ok := db[cid]
	if !ok {
		return c.JSON(http.StatusNotFound, createError())
	}
	if err := c.Bind(oldUser); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, newCustomerResponse(oldUser))
}

func getAll(c echo.Context) error {
	ans := newCustomerList(db)
	if ans.Size == 0 {
		return c.JSON(http.StatusNotFound, createError())
	}
	return c.JSON(http.StatusOK, ans)
}

func DeleteCustomer(c echo.Context) error {
	cidI, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		return err
	}
	cid := uint64(cidI)

	_, ok := db[cid]
	if !ok {
		return c.JSON(http.StatusNotFound, createError())
	}
	delete(db, cid)
	return c.JSON(http.StatusOK, JustMSG{Message: "success"})

}

func sameMonth(m int, c *Customer) bool {
	return int(c.RegisterDate.Month()) == m
}

func ReportByMonth(c echo.Context) error {
	monthI, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, createError())
	}
	var ans uint64 = 0
	for _, customer := range db {
		if sameMonth(monthI, customer) {
			ans++
		}
	}
	if ans == 0 {
		return c.JSON(http.StatusNotFound, createError())
	} else {
		return c.JSON(http.StatusOK, newReportResponse(ans))
	}

}

func searchCustomer(c echo.Context) error {
	prefix := c.QueryParam("cName")
	for _, customer := range db {
		if strings.HasPrefix(customer.Name, prefix) {
			return c.JSON(http.StatusFound, newCustomerResponse(customer))
		}
	}
	return c.JSON(http.StatusNotFound, createError())
}

func GetCustomers(c echo.Context) error {
	qp := c.QueryParam("cName")
	ok := qp != ""
	if ok {
		return searchCustomer(c)
	} else {
		return getAll(c)
	}

}

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/customers", CreateNewCustomer)
	e.PUT("/customers/:cid", EditCustomer)
	e.GET("/customers", GetCustomers)
	e.DELETE("/customers/:cid", DeleteCustomer)
	e.GET("/report/:month", ReportByMonth)

	e.Logger.Fatal(e.Start(":1323"))
}
