package models

import (
	"gorm.io/gorm"
	"justus/global"
)

type Like struct {
	Model
	Uid int `json:"uid" gorm:"column:uid"`
	Pid int `json:"p_id" gorm:"column:p_id"`
}
type LikeFormated struct {
	ID        int    `json:"id"`
	Uid       int    `json:"uid"`
	Pid       int    `json:"p_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	CreatedAt int    `json:"created_at"`
}

func (l *Like) Format() *LikeFormated {
	if l.ID <= 0 {
		return nil
	}

	return &LikeFormated{
		ID:  l.ID,
		Uid: l.Uid,
		Pid: l.Pid,
	}

}

func (l *Like) GetLikeNumber() (int, error) {
	var like Like
	var count int64
	if l.Pid > 0 {
		err := db.Model(&like).Where("p_id = ? ", l.Pid).Count(&count).Error
		if err != nil {
			global.Logger.Errorf("getCommentNumber error: %v", err)
			return 0, nil
		}
	} else {
		return 0, nil
	}
	return int(count), nil
}

func (l *Like) AddLike() (*Like, error) {
	var like Like
	if l.Pid > 0 && l.Uid > 0 {
		err := db.Model(&like).Create(l).Error
		if err != nil {
			global.Logger.Errorf("addLike error: %v", err)
			return nil, nil
		}
	} else {
		return nil, nil
	}
	return &like, nil
}

func (l *Like) GetLikeInfo() (*Like, error) {
	var like Like
	if l.Pid > 0 && l.Uid > 0 {
		err := db.Model(&like).Where("p_id = ? AND uid = ?", l.Pid, l.Uid).Last(&like).Error
		if err != nil {
			global.Logger.Errorf("getLikeInfo error: %v", err)
			return nil, nil
		}
	} else {
		return nil, nil
	}

	return &like, nil
}

// 获取图片的点赞列表
func (l *Like) GetLikeList(piDs []int) ([]*Like, error) {
	var like []*Like
	if len(piDs) > 0 {
		for _, v := range piDs {
			//fmt.Println(v)
			var likes []*Like
			err := db.Model(&like).Where("p_id = ?", v).Find(&likes).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				global.Logger.Errorf("getLikeList error: %v", err)
				return nil, nil
			}
			like = append(like, likes...)
		}

	} else {
		return nil, nil
	}

	return like, nil
}

// 获取好友的点赞列表
func (l *Like) GetFriendLikeList(piDs []int, friendIDs []int) ([]*Like, error) {
	var like []*Like
	if len(piDs) > 0 {
		for _, v := range piDs {
			//fmt.Println(v)
			var likes []*Like
			err := db.Model(&like).Where("p_id = ?", v).Where("uid in (?)", friendIDs).Find(&likes).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				global.Logger.Errorf("getLikeList error: %v", err)
				return nil, nil
			}
			like = append(like, likes...)
		}

	} else {
		return nil, nil
	}

	return like, nil
}

func (l *Like) CancelLike() (bool, error) {
	if l.Pid > 0 && l.Uid > 0 {
		err := db.Model(&Like{}).Where("p_id = ? AND uid = ?", l.Pid, l.Uid).Delete(&Like{}).Error
		if err != nil {
			global.Logger.Errorf("cancelLike error: %v", err)
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

func (l *Like) GetBatchLikeInfo(piDs []int) []*Like {
	var likes []*Like
	if l.Uid > 0 && len(piDs) > 0 {
		err := db.Model(&Like{}).Where("uid = ? AND p_id IN (?)", l.Uid, piDs).Find(&likes).Error
		if err != nil {
			global.Logger.Errorf("getBatchLikeInfo error: %v", err)
			return nil
		}
		return likes
	} else {
		return nil
	}
}
