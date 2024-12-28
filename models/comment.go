package models

import "justus/global"

type Comment struct {
	Model
	Uid     int    `json:"uid" gorm:"column:uid"`
	PId     int    `json:"p_id" gorm:"column:p_id"`
	Content string `json:"content" gorm:"column:content"`
	IsTop   int    `json:"is_top" gorm:"column:is_top"`
}
type CommentFormated struct {
	ID        int    `json:"id"`
	Uid       int    `json:"uid"`
	Pid       int    `json:"p_id"`
	Content   string `json:"content"`
	IsTop     int    `json:"is_top"`
	IsDel     int    `json:"is_del"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	CreatedAt int    `json:"created_at"`
}
type QueryCommentFirst struct {
	PId int    `json:"p_id"`
	Id  string `json:"id"`
}

func (m *Comment) Format() *CommentFormated {
	if m.ID <= 0 {
		return nil
	}

	return &CommentFormated{
		ID:        m.ID,
		Uid:       m.Uid,
		Pid:       m.PId,
		Content:   m.Content,
		IsTop:     m.IsTop,
		CreatedAt: m.CreatedAt,
	}
}

func (c *Comment) GetCommentNumber() (int, error) {
	var comment Comment
	var count int64
	if c.PId > 0 {
		err := db.Model(&comment).Where("p_id = ? ", c.PId).Count(&count).Error
		if err != nil {
			global.Logger.Errorf("getCommentNumber error: %v", err)
			return 0, nil
		}
	} else {
		return 0, nil
	}
	return int(count), nil
}

func (c *Comment) GetLastComment() (string, error) {
	var comment Comment
	if c.PId > 0 {
		err := db.Where("p_id = ? ", c.PId).Order("id desc").First(&comment).Error
		if err != nil {
			global.Logger.Errorf("GetLastComment error: %v", err)
			return "", nil
		}
	} else {
		return "", nil
	}
	return comment.Content, nil
}

func (c *Comment) CreatComment() (*Comment, error) {
	err := db.Create(&c).Error
	if err != nil {
		global.Logger.Errorf("CreatComment error: %v", err)
		return c, nil
	}
	return c, nil
}

func (c *Comment) DelComment() (bool, error) {
	err := db.Model(&c).Where("id = ? and p_id = ? and uid = ?", c.ID, c.PId, c.Uid).Delete(c).Error
	if err != nil {
		global.Logger.Errorf("DelComment error: %v", err)
		return false, nil
	}
	return true, nil
}

func (c *Comment) GetComment() (*Comment, error) {
	err := db.Model(&c).Where("id = ?", c.ID).First(&c).Error
	if err != nil {
		global.Logger.Errorf("GetComment error: %v", err)
		return c, nil
	}
	return c, nil

}

// GetCommentList 获取评论列表
func (c *Comment) GetCommentList(page int, pageSize int) ([]*Comment, error) {
	var comments []*Comment
	if c.PId > 0 {
		err := db.Where("p_id = ?", c.PId).Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments).Error
		if err != nil {
			global.Logger.Errorf("GetCommentList error: %v", err)
			return comments, nil
		}
	} else {
		return comments, nil
	}
	return comments, nil
}

// GetCommentFirst 批量获取首条评论
func (c *Comment) GetCommentFirst(ids []int) ([]*QueryCommentFirst, error) {
	var comments []*QueryCommentFirst
	err := db.Raw("SELECT a.p_id,MAX(a.t_id) AS id FROM (SELECT id,p_id,is_top,concat(is_top,',',id) AS t_id FROM yq_comment WHERE p_id IN (?)) AS a GROUP BY a.p_id", ids).Scan(&comments).Error
	if err != nil {
		global.Logger.Errorf("GetCommentFirst error: %v", err)
		return comments, err
	}
	return comments, nil
}

// GetCommentFirst 批量获取好友首条评论
func (c *Comment) GetCommentFriendFirst(ids []int, friendIDs []int) ([]*QueryCommentFirst, error) {
	var comments []*QueryCommentFirst
	err := db.Raw("SELECT a.p_id,MAX(a.t_id) AS id FROM (SELECT id,p_id,uid,is_top,concat(is_top,',',id) AS t_id FROM yq_comment WHERE p_id IN (?) AND  uid IN(?)) AS a GROUP BY a.p_id", ids, friendIDs).Scan(&comments).Error
	if err != nil {
		global.Logger.Errorf("GetCommentFirst error: %v", err)
		return comments, err
	}
	return comments, nil
}

// 批量获取图片的评论
func (c *Comment) GetCommentListByPId(pIds []int) ([]*Comment, error) {
	var comment []*Comment
	err := db.Model(&comment).Where("id IN (?)", pIds).Find(&comment).Error
	if err != nil {
		global.Logger.Errorf("GetCommentListByPId error: %v", err)
		return comment, err
	}
	return comment, nil
}
