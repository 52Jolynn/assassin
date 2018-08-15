package model

//俱乐部表
type Club struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Remark     *string `json:"remark"`
	Address    *string `json:"address"`
	Tel        *string `json:"tel"`
	CreateTime string  `json:"create_time"`
	Status     string  `json:"status"` //状态, N: 正常, D: 禁用
}

//俱乐部场地表
type Ground struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Remark     *string `json:"remark"`
	Ttype      string  `json:"ttype"`
	ClubId     int     `json:"club_id"`
	CreateTime string  `json:"create_time"`
	Status     string  `json:"status"` //N: 正常,  D: 禁用
}

//场地租用表
type GroundRental struct {
	Id          int64  `json:"id"`
	FromTime    string `json:"from_time"`
	ToTime      string `json:"to_time"`
	GroundId    int    `json:"ground_id"`
	ClubId      int    `json:"club_id"`
	RentAmount  int    `json:"rent_amount"`
	RelRentalId int64  `json:"rel_rental_id"`
	CreateTime  string `json:"create_time"`
	Status      string `json:"status"` //状态, N: 正常, L: 锁定, R:已出租
}

//球队表
type Team struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	Remark          *string `json:"remark"`
	CaptainName     *string `json:"captain_name"`
	CaptainMobile   *string `json:"captain_mobile"`
	ManagerUsername *string `json:"manager_username"`
	ManagerPasswd   *string `json:"manager_passwd"`
	CreateTime      string  `json:"create_time"`
	Status          string  `json:"status"` //状态, N: 正常, D: 禁用
}

//球队与俱乐部表
type TeamOfClub struct {
	Id             int    `json:"id"`
	TeamId         int    `json:"team_id"`
	ClubId         int    `json:"club_id"`
	PresentBalance int64  `json:"present_balance"`
	UsedAmount     int64  `json:"used_amount"`
	JoinTime       string `json:"join_time"`
	CreateTime     string `json:"create_time"`
}

//球队球衣表
type JerseyOfTeam struct {
	Id         int     `json:"id"`
	TeamId     int     `json:"team_id"`
	HomeColor  string  `json:"home_color"`
	AwayColor  *string `json:"away_color"`
	CreateTime string  `json:"create_time"`
	Status     string  `json:"status"`
}

//优惠券表
type Coupon struct {
	Id            int     `json:"id"`
	TeamId        int     `json:"team_id"`
	ClubId        int     `json:"club_id"`
	Amount        int64   `json:"amount"`
	LeastAmount   int64   `json:"least_amount"` //起用金额
	EffectiveTime string  `json:"effective_time"`
	UsedTime      *string `json:"used_time"`
	ExpireTime    *string `json:"expire_time"`
	CreateTime    string  `json:"create_time"`
	Status        string  `json:"status"` //状态, N: 正常, E: 已过期, U: 已使用
}

//球员表
type Player struct {
	Id          int     `json:"id"`
	Username    *string `json:"username"`
	Passwd      *string `json:"passwd"`
	WxOpenId    *string `json:"wx_open_id"`
	Name        string  `json:"name"`
	Remark      *string `json:"remark"`
	Mobile      *string `json:"mobile"`
	Pos         *string `json:"pos"`
	Height      float32 `json:"height"`
	Age         int     `json:"age"`
	PassVal     float32 `json:"pass_val"`
	ShotVal     float32 `json:"shot_val"`
	StrengthVal float32 `json:"strength_val"`
	DribbleVal  float32 `json:"dribble_val"`
	SpeedVal    float32 `json:"speed_val"`
	TackleVal   float32 `json:"tackle_val"`
	HeadVal     float32 `json:"head_val"`
	ThrowingVal float32 `json:"throwing_val"`
	ReactionVal float32 `json:"reaction_val"`
	CreateTime  string  `json:"create_time"`
	Status      string  `json:"status"` //状态, N: 正常, E: 退出, D: 禁用
	Level       string  `json:"level"`  //N: 普通队员(只能查看个人相关数据), S: 正式队员(可查看球队相关数据)
}

//球员数值评估表(自评+他评)
type PlayerEvaluation struct {
	Id               int     `json:"id"`
	PlayerId         int     `json:"player_id"`
	TeamId           int     `json:"team_id"`
	EvaluatePlayerId int     `json:"evaluate_player_id"`
	Fit              float32 `json:"fit"`
	CreateTime       string  `json:"create_time"`
}

//球员与球队表
type PlayerOfTeam struct {
	Id             int    `json:"id"`
	PlayerId       int    `json:"player_id"`
	TeamId         int    `json:"team_id"`
	No             string `json:"no"`
	PresentBalance int64  `json:"present_balance"`
	UsedAmount     int64  `json:"used_amount"`
	JoinTime       string `json:"join_time"`
	CreateTime     string `json:"create_time"`
}

//比赛表
type GameOfMatch struct {
	Id              int64  `json:"id"`
	HomeTeamId      int    `json:"home_team_id"`
	AwayTeamId      int    `json:"away_team_id"`
	ClubId          int    `json:"club_id"`
	GroundId        int    `json:"ground_id"`
	HomeJerseyColor string `json:"home_jersey_color"`
	AwayJerseyColor string `json:"away_jersey_color"`
	OpenTime        string `json:"open_time"`
	EnrollStartTime string `json:"enroll_start_time"`
	EnrollEndTime   string `json:"enroll_end_time"`
	EnrollQuota     int    `json:"enroll_quota"`
	RentCost        int64  `json:"rent_cost"`
	MatchDuration   int    `json:"match_duration"`
	CreateTime      string `json:"create_time"`
	Status          string `json:"status"` //状态, O: 未开赛, C: 取消, P: 开赛进行中, E: 已结束
}

//比赛报名表
type EnrollOfMatch struct {
	Id              int    `json:"id"`
	MatchId         int64  `json:"match_id"`
	PlayerId        int    `json:"player_id"`
	TemporaryPlayer int    `json:"temporary_player"` //携带散兵数
	CreateTime      string `json:"create_time"`
	Status          string `json:"status"` //报名状态, F: 报名失败, S: 报名成功, C: 取消报名
}

//球队比赛统计表
type TeamStatOfMatch struct {
	Id          int64  `json:"id"`
	MatchId     int64  `json:"match_id"`
	Ttype       string `json:"ttype"` //主队客队, home or away
	TeamId      int    `json:"team_id"`
	Score       int    `json:"score"`
	RentAmount  int64  `json:"rent_amount"`
	Shot        int    `json:"shot"`
	Foul        int    `json:"foul"`
	FreeKick    int    `json:"free_kick"`
	PenaltyKick int    `json:"penalty_kick"`
	Offside     int    `json:"offside"`
	Corner      int    `json:"corner"`
	Pass        int    `json:"pass"`
	YellowCard  int    `json:"yellow_card"`
	RedCard     int    `json:"red_card"`
	CreateTime  string `json:"create_time"`
}

//球员比赛统计表
type PlayerStatOfMatch struct {
	Id                        int64  `json:"id"`
	MatchId                   int64  `json:"match_id"`
	PlayerId                  int    `json:"player_id"`
	RentAmount                int64  `json:"rent_amount"`
	TemporaryPlayerRentAmount int64  `json:"temporary_player_rent_amount"`
	Score                     int    `json:"score"`
	Shot                      int    `json:"shot"`
	Assists                   int    `json:"assists"`
	Foul                      int    `json:"foul"`
	BreakThrough              int    `json:"break_through"`
	Tackle                    int    `json:"tackle"`
	RellowCard                int    `json:"yellow_card"`
	RedCard                   int    `json:"red_card"`
	CreateTime                string `json:"create_time"`
	PlayerStatus              string `json:"player_status"` //到场状态, N: 到场, X: 未到场
	PayBySb                   string `json:"pay_by_sb"`     //是否由他人代交场租, Y or N
	PayPlayerId               *int   `json:"pay_player_id"` //代付场租球员id
}

//球队账目表
type AccountingOfTeam struct {
	Id            string `json:"id"`
	RefId         string `json:"ref_id"`
	MatchId       int64  `json:"match_id"`
	TeamId        int    `json:"team_id"`
	Amount        int64  `json:"amount"`
	Remark        string `json:"remark"`
	BeforeBalance int64  `json:"before_balance"`
	AfterBalance  int64  `json:"after_balance"`
	Ttype         string `json:"ttype"` //记账类型, SR(收入), ZZ(支出), CZ(冲正)
	BillDate      string `json:"bill_date"`
	CreateTime    string `json:"create_time"`
}

//球员账目表
type AccountingOfPlayer struct {
	Id            string `json:"id"`
	RefId         string `json:"ref_id"`
	MatchId       int64  `json:"match_id"`
	TeamId        int    `json:"team_id"`
	PlayerId      int    `json:"player_id"`
	Amount        int64  `json:"amount"`
	Remark        string `json:"remark"`
	BeforeBalance int64  `json:"before_balance"`
	AfterBalance  int64  `json:"after_balance"`
	Ttype         string `json:"ttype"` //记账类型, R(充值), C(消费), W(提现), CZ(冲正)
	BillDate      string `json:"bill_date"`
	CreateTime    string `json:"create_time"`
}
