package main

type ItemRepositoryAbs interface{}

type ItemServiceAbs interface {
	createItem(item Item) (Item, error)
	getItem(id int) (Item, error)
	getItems() ([]Item, error)
	updateItem(id int, newItem Item) (Item, error)
	deleteItem(id int) error
}

type ItemService struct {
	itemRepository ItemRepositoryAbs
}

func (itemService *ItemService) createItem(item Item) (Item, error) {
	if err := DB.Create(&item).Error; err != nil {
		return Item{}, err
	}
	return item, nil
}

func (itemService *ItemService) getItem(id int) (Item, error) {
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		return Item{}, err
	}
	return item, nil
}

func (itemService *ItemService) getItems() ([]Item, error) {
	var items []Item
	if err := DB.Find(&items).Error; err != nil {
		return []Item{}, err
	}
	return items, nil
}

func (itemService *ItemService) updateItem(id int, newItem Item) (Item, error) {
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		return Item{}, err
	}
	DB.Model(&item).Updates(newItem)
	return item, nil
}

func (itemService *ItemService) deleteItem(id int) error {
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		return err
	}
	DB.Delete(item)
	return nil
}
