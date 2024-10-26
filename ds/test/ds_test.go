package ds_test

import (
	"testing"
	"time"

	"github.com/QingShan-Xu/web/ds"
)

// Address 结构体，用于嵌套字段测试
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

// Profile 结构体，用于切片字段测试
type Profile struct {
	Bio string `json:"bio"`
}

// User 结构体，包含嵌套字段和切片字段
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   Address   `json:"address"`
	Profiles  []Profile `json:"profiles"`
	CreatedAt time.Time `json:"created_at"`
}

// Employee 结构体，嵌入 User 结构体，用于匿名字段测试
type Employee struct {
	User            // 匿名嵌入字段
	Position string `json:"position"`
}

// Product 结构体，用于多标签测试
type Product struct {
	ID    int     `json:"id" db:"product_id"`
	Name  string  `json:"name" db:"product_name"`
	Price float64 `json:"price" db:"product_price"`
}

// TestStructReader_Advanced 测试 StructReader 的多种功能
func TestStructReader_Advanced(t *testing.T) {
	// 初始化 Employee 实例，用于测试嵌入字段和嵌套字段
	employee := Employee{
		User: User{
			ID:   1,
			Name: "John Doe",
			Address: Address{
				Street: "123 Main St",
				City:   "Metropolis",
			},
			Profiles: []Profile{
				{Bio: "Developer"},
				{Bio: "Blogger"},
			},
			CreatedAt: time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		},
		Position: "Software Engineer",
	}

	// 创建 StructReader 实例
	reader, err := ds.NewStructReader(employee)
	if err != nil {
		t.Fatalf("Failed to create StructReader: %v", err)
	}

	// 子测试 1：访问顶层字段多次
	t.Run("TopLevelFields", func(t *testing.T) {
		// 第一次访问 Name 字段
		nameField, err := reader.GetField("Name")
		if err != nil {
			t.Fatalf("Failed to get field 'Name': %v", err)
		}
		name, ok := nameField.SafeString()
		if !ok || name != "John Doe" {
			t.Errorf("Expected Name 'John Doe', got '%v'", name)
		}
		t.Logf("Name: %s", name)

		// 第二次访问 Name 字段
		nameField2, err := reader.GetField("Name")
		if err != nil {
			t.Fatalf("Failed to get field 'Name' again: %v", err)
		}
		name2, ok := nameField2.SafeString()
		if !ok || name2 != "John Doe" {
			t.Errorf("Expected Name 'John Doe', got '%v'", name2)
		}
		t.Logf("Name (second access): %s", name2)
	})

	// 子测试 2：访问嵌套字段多次
	t.Run("NestedFields", func(t *testing.T) {
		// 第一次访问 Address.Street 字段
		streetField, err := reader.GetField("Address.Street")
		if err != nil {
			t.Fatalf("Failed to get field 'Address.Street': %v", err)
		}
		street, ok := streetField.SafeString()
		if !ok || street != "123 Main St" {
			t.Errorf("Expected Street '123 Main St', got '%v'", street)
		}
		t.Logf("Street: %s", street)

		// 第二次访问 Address.Street 字段
		streetField2, err := reader.GetField("Address.Street")
		if err != nil {
			t.Fatalf("Failed to get field 'Address.Street' again: %v", err)
		}
		street2, ok := streetField2.SafeString()
		if !ok || street2 != "123 Main St" {
			t.Errorf("Expected Street '123 Main St', got '%v'", street2)
		}
		t.Logf("Street (second access): %s", street2)
	})

	// 子测试 3：使用 ._len 语法获取切片长度
	t.Run("LengthAccess", func(t *testing.T) {
		// 获取 Profiles 切片的长度
		profilesLenField, err := reader.GetField("Profiles._len")
		if err != nil {
			t.Fatalf("Failed to get field 'Profiles._len': %v", err)
		}
		profilesLen, ok := profilesLenField.SafeInt()
		if !ok || profilesLen != 2 {
			t.Errorf("Expected Profiles length 2, got '%v'", profilesLen)
		}
		t.Logf("Profiles length: %d", profilesLen)

		// 再次获取 Profiles 切片的长度
		profilesLenField2, err := reader.GetField("Profiles._len")
		if err != nil {
			t.Fatalf("Failed to get field 'Profiles._len' again: %v", err)
		}
		profilesLen2, ok := profilesLenField2.SafeInt()
		if !ok || profilesLen2 != 2 {
			t.Errorf("Expected Profiles length 2, got '%v'", profilesLen2)
		}
		t.Logf("Profiles length (second access): %d", profilesLen2)
	})

	// 子测试 4：数组/切片映射语法
	t.Run("ArrayMapping", func(t *testing.T) {
		// 使用数组映射语法提取 Profiles 切片中的 Bio 字段
		mappedProfilesField, err := reader.GetField("Profiles[Bio]")
		if err != nil {
			t.Fatalf("Failed to get field 'Profiles[Bio]': %v", err)
		}
		mappedProfiles, ok := mappedProfilesField.Interface().([]map[string]interface{})
		if !ok || len(mappedProfiles) != 2 {
			t.Errorf("Expected mapped Profiles of length 2, got '%v'", mappedProfiles)
		}
		t.Logf("Mapped Profiles: %+v", mappedProfiles)
		t.Log(mappedProfiles)

		// 再次使用数组映射语法提取 Profiles 切片中的 Bio 字段
		mappedProfilesField2, err := reader.GetField("Profiles[Bio]")
		if err != nil {
			t.Fatalf("Failed to get field 'Profiles[Bio]' again: %v", err)
		}
		mappedProfiles2, ok := mappedProfilesField2.Interface().([]map[string]interface{})
		if !ok || len(mappedProfiles2) != 2 {
			t.Errorf("Expected mapped Profiles of length 2, got '%v'", mappedProfiles2)
		}
		t.Logf("Mapped Profiles (second access): %+v", mappedProfiles2)
	})

	// 子测试 5：访问嵌入（匿名）字段
	t.Run("EmbeddedFields", func(t *testing.T) {
		// 访问 Position 字段
		positionField, err := reader.GetField("Position")
		if err != nil {
			t.Fatalf("Failed to get field 'Position': %v", err)
		}
		position, ok := positionField.SafeString()
		if !ok || position != "Software Engineer" {
			t.Errorf("Expected Position 'Software Engineer', got '%v'", position)
		}
		t.Logf("Position: %s", position)

		// 访问嵌入字段 Name
		nameField, err := reader.GetField("Name")
		if err != nil {
			t.Fatalf("Failed to get embedded field 'Name': %v", err)
		}
		name, ok := nameField.SafeString()
		if !ok || name != "John Doe" {
			t.Errorf("Expected Name 'John Doe', got '%v'", name)
		}
		t.Logf("Embedded Name: %s", name)

		// 访问嵌入字段 Address.City
		cityField, err := reader.GetField("Address.City")
		if err != nil {
			t.Fatalf("Failed to get embedded field 'Address.City': %v", err)
		}
		city, ok := cityField.SafeString()
		if !ok || city != "Metropolis" {
			t.Errorf("Expected City 'Metropolis', got '%v'", city)
		}
		t.Logf("Embedded Address City: %s", city)
	})

	// 子测试 6：访问 time.Time 字段
	t.Run("TimeField", func(t *testing.T) {
		createdAtField, err := reader.GetField("CreatedAt")
		if err != nil {
			t.Fatalf("Failed to get field 'CreatedAt': %v", err)
		}
		createdAt, ok := createdAtField.SafeTime()
		if !ok {
			t.Errorf("Failed to get CreatedAt as time.Time")
		}
		t.Logf("CreatedAt: %s", createdAt.String())
	})

	// 子测试 7：访问不存在的字段
	t.Run("NonExistentField", func(t *testing.T) {
		_, err := reader.GetField("NonExistent")
		if err == nil {
			t.Errorf("Expected error when accessing non-existent field, but got none")
		} else {
			t.Logf("Correctly received error: %v", err)
		}
	})

	// 子测试 8：访问具有多个标签的字段
	t.Run("MultipleTags", func(t *testing.T) {
		// 初始化 Product 实例
		product := Product{
			ID:    101,
			Name:  "Laptop",
			Price: 999.99,
		}

		// 创建 StructReader 实例
		productReader, err := ds.NewStructReader(product)
		if err != nil {
			t.Fatalf("Failed to create StructReader for Product: %v", err)
		}

		// 访问 ID 字段并验证多个标签
		idField, err := productReader.GetField("ID")
		if err != nil {
			t.Fatalf("Failed to get field 'ID': %v", err)
		}
		tags := idField.GetTag()
		if jsonTag, ok := tags["json"]; !ok || jsonTag.Value != "id" {
			t.Errorf("Expected json tag 'id', got '%v'", jsonTag)
		}
		if dbTag, ok := tags["db"]; !ok || dbTag.Value != "product_id" {
			t.Errorf("Expected db tag 'product_id', got '%v'", dbTag)
		}
		t.Logf("Product ID tags: %+v", tags)

		// 访问 Name 字段并验证多个标签
		nameField, err := productReader.GetField("Name")
		if err != nil {
			t.Fatalf("Failed to get field 'Name': %v", err)
		}
		nameTags := nameField.GetTag()
		if jsonTag, ok := nameTags["json"]; !ok || jsonTag.Value != "name" {
			t.Errorf("Expected json tag 'name', got '%v'", jsonTag)
		}
		if dbTag, ok := nameTags["db"]; !ok || dbTag.Value != "product_name" {
			t.Errorf("Expected db tag 'product_name', got '%v'", dbTag)
		}
		t.Logf("Product Name tags: %+v", nameTags)

		// 访问 Price 字段并验证多个标签
		priceField, err := productReader.GetField("Price")
		if err != nil {
			t.Fatalf("Failed to get field 'Price': %v", err)
		}
		priceTags := priceField.GetTag()
		if jsonTag, ok := priceTags["json"]; !ok || jsonTag.Value != "price" {
			t.Errorf("Expected json tag 'price', got '%v'", jsonTag)
		}
		if dbTag, ok := priceTags["db"]; !ok || dbTag.Value != "product_price" {
			t.Errorf("Expected db tag 'product_price', got '%v'", dbTag)
		}
		t.Logf("Product Price tags: %+v", priceTags)
	})

	// 子测试 9：多次访问嵌入字段
	t.Run("EmbeddedFieldsMultipleAccess", func(t *testing.T) {
		// 多次访问嵌入字段的 Name 和 Address.City
		for i := 1; i <= 3; i++ {
			nameField, err := reader.GetField("Name")
			if err != nil {
				t.Fatalf("Failed to get embedded field 'Name' on iteration %d: %v", i, err)
			}
			name, ok := nameField.SafeString()
			if !ok || name != "John Doe" {
				t.Errorf("Expected Name 'John Doe', got '%v' on iteration %d", name, i)
			}
			t.Logf("Iteration %d: Embedded Name: %s", i, name)

			cityField, err := reader.GetField("Address.City")
			if err != nil {
				t.Fatalf("Failed to get embedded field 'Address.City' on iteration %d: %v", i, err)
			}
			city, ok := cityField.SafeString()
			if !ok || city != "Metropolis" {
				t.Errorf("Expected City 'Metropolis', got '%v' on iteration %d", city, i)
			}
			t.Logf("Iteration %d: Embedded Address City: %s", i, city)
		}
	})
}
