package ranking

type Character struct {
    CharID    int    `json:"char_id"`
    AccID     int    `json:"acc_id"`
    ClassID   int    `json:"class_id"`
    ClassName string `json:"class_name"`
}

type Score struct {
    ScoreID      int    `json:"score_id"`
    CharID       int    `json:"char_id"`
    RewardScore  int    `json:"reward_score"`
    Username     string `json:"username"`
    ClassName    string `json:"class_name"`
}

type RankingEntry struct {
    Rank         int    `json:"rank"`
    Username     string `json:"username"`
    ClassName    string `json:"class_name"`
    RewardScore  int    `json:"reward_score"`
}

type Class struct {
    ID          int    `json:"id"`
    RaceID      int    `json:"race_id"`
    RaceName    string `json:"race_name"`
    Name        string `json:"name"`
    Title       string `json:"title"`
    Description string `json:"description"`
    CombatType  string `json:"combat_type"`
    Damage      int    `json:"damage"`
    Defense     int    `json:"defense"`
    Difficulty  int    `json:"difficulty"`
    Speed       int    `json:"speed"`
}

type RankingResponse struct {
    Rankings    []RankingEntry `json:"rankings"`
    TotalCount  int           `json:"total_count"`
    CurrentPage int           `json:"current_page"`
    TotalPages  int           `json:"total_pages"`
}
