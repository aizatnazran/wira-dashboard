package ranking

import (
    "database/sql"
    "fmt"
    "context"
    "wira-assignment/cache"
    "time"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) GetRankings(classID, page, limit int, search string) ([]RankingEntry, int, error) {
    ctx := context.Background()
    
    // Create cache key based on parameters
    cacheKey := fmt.Sprintf("rankings:%d:%d:%d:%s", classID, page, limit, search)
    
    // Try to get from cache
    var cachedResult struct {
        Rankings []RankingEntry
        Total    int
    }
    err := cache.Get(ctx, cacheKey, &cachedResult)
    if err == nil {
        return cachedResult.Rankings, cachedResult.Total, nil
    }

    offset := (page - 1) * limit

    // Base query
    query := `
        WITH RankedScores AS (
            SELECT 
                ROW_NUMBER() OVER (ORDER BY s.reward_score DESC) as rank,
                u.username,
                c.name as class_name,
                s.reward_score,
                COUNT(*) OVER() as total_count
            FROM scores s
            JOIN characters ch ON s.char_id = ch.char_id
            JOIN accounts u ON ch.acc_id = u.acc_id
            JOIN classes c ON ch.class_id = c.id
            WHERE 1=1
    `

    // Add class filter if classID is provided
    var args []interface{}
    argCount := 1
    if classID > 0 {
        query += fmt.Sprintf(" AND ch.class_id = $%d", argCount)
        args = append(args, classID)
        argCount++
    }

    // Add search condition if search is provided
    if search != "" {
        query += fmt.Sprintf(" AND LOWER(u.username) LIKE LOWER($%d)", argCount)
        args = append(args, "%"+search+"%")
        argCount++
    }

    query += `
        )
        SELECT rank, username, class_name, reward_score, total_count
        FROM RankedScores
        LIMIT $` + fmt.Sprint(argCount) + ` OFFSET $` + fmt.Sprint(argCount+1) + `
    `
    args = append(args, limit, offset)

    // Execute query
    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("error querying rankings: %v", err)
    }
    defer rows.Close()

    var rankings []RankingEntry
    var totalCount int
    for rows.Next() {
        var entry RankingEntry
        err := rows.Scan(&entry.Rank, &entry.Username, &entry.ClassName, &entry.RewardScore, &totalCount)
        if err != nil {
            return nil, 0, fmt.Errorf("error scanning ranking entry: %v", err)
        }
        rankings = append(rankings, entry)
    }

    // Cache the results for 1 minute
    result := struct {
        Rankings []RankingEntry
        Total    int
    }{
        Rankings: rankings,
        Total:    totalCount,
    }
    err = cache.Set(ctx, cacheKey, result, time.Minute)
    if err != nil {
        fmt.Printf("Warning: Failed to cache rankings: %v\n", err)
    }

    return rankings, totalCount, nil
}

func (r *Repository) GetClassIDByName(className string) (int, error) {
    var id int
    err := r.db.QueryRow("SELECT id FROM classes WHERE name = $1", className).Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("class not found: %s", className)
        }
        return 0, fmt.Errorf("error getting class ID: %v", err)
    }
    return id, nil
}

func (r *Repository) GetClasses() ([]Class, error) {
    ctx := context.Background()
    
    // Try to get from cache
    var classes []Class
    err := cache.Get(ctx, "classes", &classes)
    if err == nil {
        return classes, nil
    }

    query := `
        SELECT c.id, c.race_id, r.name as race_name, c.name, c.title, 
               c.description, c.combat_type, c.damage, c.defense, 
               c.difficulty, c.speed
        FROM classes c
        JOIN races r ON c.race_id = r.id
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error querying classes: %v", err)
    }
    defer rows.Close()

    classes = []Class{}
    for rows.Next() {
        var class Class
        err := rows.Scan(
            &class.ID, &class.RaceID, &class.RaceName, &class.Name, &class.Title,
            &class.Description, &class.CombatType, &class.Damage, &class.Defense,
            &class.Difficulty, &class.Speed,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning class: %v", err)
        }
        classes = append(classes, class)
    }

    // Cache the result for 1 hour
    err = cache.Set(ctx, "classes", classes, time.Hour)
    if err != nil {
        fmt.Printf("Warning: Failed to cache classes: %v\n", err)
    }

    return classes, nil
}

func (r *Repository) GetUserCharacters(userID int) ([]Character, error) {
    query := `
        SELECT c.char_id, c.acc_id, c.class_id, cl.name as class_name
        FROM characters c
        JOIN classes cl ON c.class_id = cl.id
        WHERE c.acc_id = $1
    `
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("error querying user characters: %v", err)
    }
    defer rows.Close()

    var characters []Character
    for rows.Next() {
        var char Character
        err := rows.Scan(&char.CharID, &char.AccID, &char.ClassID, &char.ClassName)
        if err != nil {
            return nil, fmt.Errorf("error scanning character: %v", err)
        }
        characters = append(characters, char)
    }

    return characters, nil
}

func (r *Repository) CreateCharacter(userID int, classID int) error {
    // Verify that the class exists
    var exists bool
    err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1)", classID).Scan(&exists)
    if err != nil {
        return fmt.Errorf("error checking class existence: %v", err)
    }
    if !exists {
        return fmt.Errorf("invalid class ID")
    }

    // Create the character
    query := `
        INSERT INTO characters (acc_id, class_id)
        VALUES ($1, $2)
    `
    _, err = r.db.Exec(query, userID, classID)
    if err != nil {
        return fmt.Errorf("error creating character: %v", err)
    }

    return nil
}

func (r *Repository) UpdateScore(charID int, score int) error {
    query := `
        INSERT INTO scores (char_id, reward_score)
        VALUES ($1, $2)
        ON CONFLICT (char_id) DO UPDATE
        SET reward_score = $2, updated_at = CURRENT_TIMESTAMP
    `
    _, err := r.db.Exec(query, charID, score)
    if err != nil {
        return fmt.Errorf("error updating score: %v", err)
    }

    return nil
}
