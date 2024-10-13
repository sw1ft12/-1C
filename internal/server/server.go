package server

import (
    "1C/internal/models"
    "1C/internal/repository"
    "github.com/gin-gonic/gin"
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
    "github.com/go-echarts/go-echarts/v2/types"
    "github.com/jackc/pgx/v5/pgxpool"
    "net/http"
    "time"
)

type Server struct {
    router *gin.Engine
    repo   *repository.Repo
}

func NewServer(conn *pgxpool.Pool) *Server {
    return &Server{
        router: gin.Default(),
        repo:   repository.NewRepo(conn),
    }
}

func (s *Server) AddDish(context *gin.Context) {
    var dish models.Dish
    err := context.BindJSON(&dish)
    if err != nil {
        context.JSON(http.StatusBadRequest, err)
        return
    }
    err = s.repo.CreateDish(context.Request.Context(), dish)
    if err != nil {
        context.JSON(http.StatusBadRequest, err)
        return
    }
    context.JSON(http.StatusCreated, dish)
}

func (s *Server) GetDishes(context *gin.Context) {
    dishes, err := s.repo.GetDishes(context.Request.Context())
    if err != nil {
        context.JSON(http.StatusBadRequest, err)
        return
    }
    context.JSON(http.StatusOK, dishes)
}

func (s *Server) FindDishByName(context *gin.Context) {
    dishName := context.Query("dish")
    if dishName == "" {
        context.JSON(http.StatusOK, "")
        return
    }
    dish, err := s.repo.FindDish(context.Request.Context(), dishName)
    if err != nil {
        context.JSON(http.StatusNotFound, err)
        return
    }
    context.JSON(http.StatusOK, dish)
}

func (s *Server) UpdateDishByName(context *gin.Context) {
    dishName := context.Query("dish")
    if dishName == "" {
        context.JSON(http.StatusNotFound, "")
        return
    }
    var dish models.Dish
    err := context.BindJSON(&dish)
    if err != nil {
        context.JSON(http.StatusBadRequest, err)
        return
    }
    dish, err = s.repo.UpdateDish(context.Request.Context(), dishName, dish)
    if err != nil {
        context.JSON(http.StatusBadRequest, err)
        return
    }
    context.JSON(http.StatusOK, dish)
}

func (s *Server) SelectDish(context *gin.Context) {
    dishName := context.Query("dish")
    if dishName == "" {
        context.String(http.StatusNotFound, "")
        return
    }
    err := s.repo.ChooseDish(context.Request.Context(), dishName)
    if err != nil {
        context.String(http.StatusBadRequest, err.Error())
        return
    }
    context.String(http.StatusOK, dishName+" selected")
}

func (s *Server) CalcCalories(context *gin.Context) {
    date1, date2 := context.Query("date1"), context.Query("date2")
    if date1 == "" || date2 == "" {
        context.String(http.StatusBadRequest, "")
        return
    }
    stats, err := s.repo.CalcCalories(context.Request.Context(), date1, date2)
    if err != nil {
        context.String(http.StatusBadRequest, err.Error())
        return
    }
    date := make([]time.Time, len(stats))
    sum := make([]opts.LineData, len(stats))
    for i := range stats {
        date[i] = stats[i].Date
        sum[i] = opts.LineData{
            Value: stats[i].SumCalories,
        }
    }
    line := charts.NewLine()
    line.SetGlobalOptions(
        charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
        charts.WithTitleOpts(opts.Title{
            Title: "Потребление калорий за период",
        }))
    line.SetXAxis(date).
        AddSeries("Калории", sum)
    line.Render(context.Writer)
}

func (s *Server) InitRoutes() {
    dishes := s.router.Group("/dishes")
    {
        dishes.POST("/new", s.AddDish)
        dishes.GET("/", s.GetDishes)
        dishes.GET("/search", s.FindDishByName)
        dishes.POST("/update", s.UpdateDishByName)
        dishes.POST("/select", s.SelectDish)
        dishes.GET("/stat", s.CalcCalories)
    }
}

func (s *Server) Run(address string) error {
    s.InitRoutes()
    err := s.router.Run(address)
    if err != nil {
        return err
    }
    return nil
}
