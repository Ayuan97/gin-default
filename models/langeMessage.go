package models

import (
	"gorm.io/gorm"
	"justus/global"
)

type LangeMessage struct {
	LId        int `gorm:"primary_key" json:"lid"`
	ErrorCode string `json:"error_code"`
	ZhHans string `json:"zh-Hans"`
	ZhHant string `json:"zh-Hant"`
	En string `json:"en"`
	Th string `json:"th"`
	Fr string `json:"fr"`
	Es string `json:"es"`
	Fil string `json:"fil"`
	Ms string `json:"ms"`
	Pt string `json:"pt"`
	Ja string `json:"ja"`
	Id string `json:"id"`
	Af string `json:"af"`
	Am string `json:"am"`
	Bg string `json:"bg"`
	Ca string `json:"ca"`
	Hr string `json:"hr"`
	Cs string `json:"cs"`
	Da string `json:"da"`
	Nl string `json:"nl"`
	Et string `json:"et"`
	Fi string `json:"fi"`
	De string `json:"de"`
	El string `json:"el"`
	He string `json:"he"`
	Hi string `json:"hi"`
	Hu string `json:"hu"`
	Is string `json:"is"`
	It string `json:"it"`
	Ko string `json:"ko"`
	Lv string `json:"lv"`
	Lt string `json:"lt"`
	Nb string `json:"nb"`
	Pl string `json:"pl"`
	Ro string `json:"ro"`
	Ru string `json:"ru"`
	Sr string `json:"sr"`
	Sw string `json:"sw"`
	Sv string `json:"sv"`
	Tr string `json:"tr"`
	Uk string `json:"uk"`
	Vi string `json:"vi"`
	Zu string `json:"zu"`
	Ar string `json:"ar"`
}



//入库
func (m *LangeMessage)GetLangeMessageList() ([]LangeMessage,error) {
	var langeMessage []LangeMessage
	err := db.Find(&langeMessage).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.Logger.Error("message GetList:",err)

	}
	return langeMessage,nil
}

