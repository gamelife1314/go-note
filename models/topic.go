package models

import "time"

type Topic struct {
	ID        uint   `gorm:"primary_key" json:"-"`
	Name      string `gorm:"type:char(24);not null;unique" json:"name"`
	ParentId  *uint  `gorm:"column:parent_id" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Articles []Article `gorm:"many2many:article_topic"`
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
		Database.Where(Topic{Name: parent}).FirstOrCreate(&parentModel)
		for _, child := range children {
			var childModel Topic
			Database.Where(Topic{Name: child, ParentId: &parentModel.ID}).FirstOrCreate(&childModel)
		}
	}
}
