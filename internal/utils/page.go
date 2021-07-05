package utils

import (
    "github.com/labstack/echo"
    "gorm.io/gorm"
    "math"
)

type Paginate struct {
    Page  int
    Limit int
}

// 进行分页
func NewPaginate(page int) *Paginate {
    return &Paginate{Page: page}
}

// 设置页码
func (p *Paginate) SetPage(page int) *Paginate {
    if page <= 0 {
        page = 1
    }
    p.Page = page
    return p
}

// 设置每页条数
func (p *Paginate) SetLimit(limit int) *Paginate {
    if limit <= 0 {
        limit = 1
    } else if limit > 100 {
        limit = 100
    }
    p.Limit = limit
    return p
}

// 进行分页
func (p *Paginate) Paginate(query *gorm.DB, repo interface{}) echo.Map {

    c := make(chan int64, 1)
    // query = query.Session(&gorm.Session{WithConditions: true})
    go func() {
        c <- p.Count(query)
    }()

    if p.Page == 0 {
        p.Page = 1
    }
    if p.Limit == 0 {
        p.Limit = 20
    }

    offset := (p.Page - 1) * p.Limit
    result := query.Offset(offset).Limit(p.Limit).Find(repo)
    count := <-c

    res := echo.Map{}
    res["total_count"] = count
    res["total_page"] = int(math.Ceil(float64(count) / float64(p.Limit)))
    res["current_page"] = p.Page
    res["page_limit"] = p.Limit
    res["page_count"] = int(result.RowsAffected)
    if res["page_count"] != 0 {
        res["items"] = repo
    } else {
        res["items"] = make([]interface{}, 0)
    }

    return res
}

// 统计总行数
func (p *Paginate) Count(query *gorm.DB) int64 {
    var count int64
    query.Count(&count)
    return count
}
