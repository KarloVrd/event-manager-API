package models

type Organization struct {
	ID   uint   `json:"id" gorm:"primary_key,auto_increment"`
	Name string `json:"name" required:"true" gorm:"unique"`
}

func GetOrganizationByName(name string) (Organization, error) {
	var organization Organization
	if err := DB.Where("name = ?", name).First(&organization).Error; err != nil {
		return organization, err
	}
	return organization, nil
}

func GetOrganizationByID(id uint) (Organization, error) {
	var organization Organization
	if err := DB.Where("id = ?", id).First(&organization).Error; err != nil {
		return organization, err
	}
	return organization, nil
}
