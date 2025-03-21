package sys

// 定义枚举类型
type CAtegory string

// 定义枚举常量
const (
	CAtegoryIndex CAtegory = ""               // 0,官网,official
	CAtegoryWT    CAtegory = "whisperingtime" // 1,枫迹,whisperingtime
	CAtegorySQ    CAtegory = "sq"             // 2,暂定为sq
)

func GetCategory(name string) (CAtegory, error) {
	switch name {
	case "":
		return CAtegoryIndex, nil
	case "whisperingtime":
		return CAtegoryWT, nil
	case "sq":
		return CAtegorySQ, nil
	default:
		return "", ErrCategory
	}
}
