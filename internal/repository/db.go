package repository

import (
    "1C/internal/models"
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
    conn *pgxpool.Pool
}

func NewRepo(conn *pgxpool.Pool) *Repo {
    return &Repo{
        conn: conn,
    }
}

func (r *Repo) CreateDish(ctx context.Context, dish models.Dish) error {
    query := `INSERT INTO Dishes (name, calorie, protein, fat, carbohydrate) VALUES(@name, @calorie, @protein, @fat, @carbohydrate)`
    _, err := r.conn.Query(ctx, query, pgx.NamedArgs{
        "name":         dish.Name,
        "calorie":      dish.Calorie,
        "protein":      dish.Protein,
        "fat":          dish.Fat,
        "carbohydrate": dish.Carbohydrate,
    })
    if err != nil {
        return err
    }
    return nil
}

func (r *Repo) GetDishes(ctx context.Context) ([]models.Dish, error) {
    query := `SELECT * FROM Dishes`
    rows, err := r.conn.Query(ctx, query)
    if err != nil {
        return []models.Dish{}, err
    }
    dishes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Dish])
    if err != nil {
        return []models.Dish{}, err
    }
    return dishes, nil
}

func (r *Repo) FindDish(ctx context.Context, name string) (models.Dish, error) {
    query := `SELECT * FROM Dishes WHERE name=@name`
    row, _ := r.conn.Query(ctx, query, pgx.NamedArgs{
        "name": name,
    })
    dishes, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Dish])
    if err != nil {
        return models.Dish{}, err
    }
    return dishes, nil
}

func (r *Repo) UpdateDish(ctx context.Context, name string, dish models.Dish) (models.Dish, error) {
    query := `UPDATE Dishes SET name=@name, calorie=@calorie, protein=@protein, fat=@fat, carbohydrate=@carbohydrate
                RETURNING name, calorie, protein, fat, carbohydrate`
    row, _ := r.conn.Query(ctx, query, pgx.NamedArgs{
        "name":         dish.Name,
        "calorie":      dish.Calorie,
        "protein":      dish.Protein,
        "fat":          dish.Fat,
        "carbohydrate": dish.Carbohydrate,
    })
    dishes, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Dish])
    if err != nil {
        return models.Dish{}, err
    }
    return dishes, nil
}

func (r *Repo) ChooseDish(ctx context.Context, name string) error {
    query := `INSERT INTO SelectedDished (name) VALUES(@name)`
    _, err := r.conn.Query(ctx, query)
    if err != nil {
        return err
    }
    return nil
}

func (r *Repo) CalcCalories(ctx context.Context, date1, date2 string) ([]models.Stat, error) {
    query := `SELECT sd.date, SUM(d.calorie) FROM Dishes d JOIN SelectedDished sd ON d.name = sd.name WHERE sd.date BETWEEN @date1 AND @date2 GROUP BY sd.date`
    rows, err := r.conn.Query(ctx, query, pgx.NamedArgs{
        "date1": date1,
        "date2": date2,
    })
    if err != nil {
        return nil, err
    }
    stats, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Stat])
    if err != nil {
        return []models.Stat{}, err
    }
    return stats, nil
}
