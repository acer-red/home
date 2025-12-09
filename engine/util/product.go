package util

// 定义枚举类型
type CAtegory string

// 定义枚举常量
const (
	CAtegoryOfficial CAtegory = "official"       // 0,官网,official
	CAtegoryWT       CAtegory = "whisperingtime" // 1,枫迹,whisperingtime
	CAtegorySQ       CAtegory = "sq"             // 2,暂定为sq
)

func GetCategory(name string) (CAtegory, bool) {
	switch name {
	case "official":
		return CAtegoryOfficial, true
	case "whisperingtime":
		return CAtegoryWT, true
	case "sq":
		return CAtegorySQ, true
	default:
		return "", false
	}
}
func (c CAtegory) GetAuthCookiePrefix() string {
	switch c {
	case CAtegoryOfficial:
		return "index:auth"
	case CAtegoryWT:
		return "wt:auth"
	case CAtegorySQ:
		return "sq:auth"
	default:
		return "default:auth"
	}
}
