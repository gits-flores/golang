package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Profile struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Photo  string    `gorm:"size:255;not null" json:"photo"`
	User    User      `json:"user"`
	UserID  uint32    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (p *Profile) Prepare() {
	p.ID = 0
	p.Photo = "-"
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Profile) Validate() error {

	if p.Photo == "" {
		return errors.New("Required Photo")
	}

	if p.UserID < 1 {
		return errors.New("Required User")
	}
	return nil
}

func (p *Profile) SaveProfile(db *gorm.DB) (*Profile, error) {

	var err error
	err = db.Debug().Create(&p).Error
	if err != nil {
		return &Profile{}, err
	}
	return p, nil
}

func (p *Profile) FindAllProfiles(db *gorm.DB) (*[]Profile, error) {
	var err error
	Profiles := []Profile{}
	err = db.Debug().Model(&Profile{}).Limit(100).Find(&Profiles).Error
	if err != nil {
		return &[]Profile{}, err
	}
	if len(Profiles) > 0 {
		for i, _ := range Profiles {
			err := db.Debug().Model(&User{}).Where("id = ?", Profiles[i].UserID).Take(&Profiles[i].User).Error
			if err != nil {
				return &[]Profile{}, err
			}
		}
	}
	return &Profiles, nil
}

func (p *Profile) FindProfileByID(db *gorm.DB, pid uint64) (*Profile, error) {
	var err error
	err = db.Debug().Model(&Profile{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Profile{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Profile{}, err
		}
	}
	return p, nil
}

func (p *Profile) UpdateAProfile(db *gorm.DB) (*Profile, error) {

	var err error

	err = db.Debug().Model(&Profile{}).Where("id = ?", p.ID).Updates(Profile{Photo: p.Photo, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Profile{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Profile{}, err
		}
	}
	return p, nil
}

func (p *Profile) DeleteAProfile(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Profile{}).Where("id = ? and user_id = ?", pid, uid).Take(&Profile{}).Delete(&Profile{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Profile not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}