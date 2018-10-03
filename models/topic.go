package models

import (
	"time"
)

type Topic struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Name      string     `gorm:"type:char(24);not null;unique" json:"name"`
	ParentId  *uint      `gorm:"column:parent_id" json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Articles []Article `gorm:"many2many:article_topic" json:"-"`
}

func InitTopics() {
	var topics = map[string][]string{
		"专题": {"EasyWeChat"},
		"社区": {"规范", "公告"},
		"微信": {"公众号", "企业微信", "开发平台", "小程序", "小游戏", "支付", "调试", "审核", "框架", "开发者工具",
			"客户端", "服务端", "教程", "API", "朋友圈", "组件", "开源推荐", "自定义菜单", "消息", "小程序码",
			"数据管理", "统计分析", "模板消息", "安全", "解密", "通讯录", "素材", "电子发票", "SDK", "缓存",
			"群发", "二维码", "卡券", "门店", "客服", "摇一摇", "微信硬件", "设备", "签名", "统计", "语义理解",
			"测试号", "商户", "UI", "WeUI", "JSSDK", "网页授权", "分享"},
		"语言":   {"GO", "Javascript", "PHP", "Python", "CSS", "HTML"},
		"技术框架": {"Iris", "Vue", "React", "Laravel"},
	}

	for parent, children := range topics {
		var parentModel Topic
		Database.Where(Topic{Name: parent, ParentId: nil}).FirstOrCreate(&parentModel)
		for _, child := range children {
			var childModel Topic
			Database.Where(Topic{Name: child, ParentId: &parentModel.ID}).FirstOrCreate(&childModel)
		}
	}
}

func TopicsByLevel() []map[string]interface{} {
	var parents []Topic
	var results []map[string]interface{}
	Database.Where(map[string]interface{}{"parent_id": nil}).Find(&parents)
	for _, parent := range parents {
		var children []Topic
		var childrenArr []map[string]interface{}
		Database.Where(map[string]interface{}{"parent_id": parent.ID}).Find(&children)
		for _, child := range children {
			childrenArr = append(childrenArr, map[string]interface{}{
				"id":   child.ID,
				"name": child.Name,
			})
		}
		results = append(results, map[string]interface{}{
			"id":       parent.ID,
			"name":     parent.Name,
			"children": childrenArr,
		})
	}
	return results
}

func (t *Topic) Transform() map[string]interface{} {
	return map[string]interface{}{
		"id":   t.ID,
		"name": t.Name,
	}
}
