// Package mapper 提供领域模型与持久化对象之间的通用转换能力
// 支持 GORM GEN / GORM / Ent 等多种 ORM 的映射
package mapper

import (
	"github.com/jinzhu/copier"
)

// Mapper 定义领域模型(Domain)与持久化对象(Entity)之间的转换接口
// D = Domain model (biz层)
// E = Entity/PO (data层)
type Mapper[D any, E any] interface {
	// ToDomain 将持久化对象转换为领域模型
	ToDomain(entity *E) *D

	// ToEntity 将领域模型转换为持久化对象
	ToEntity(domain *D) *E

	// ToDomainList 批量转换持久化对象为领域模型
	ToDomainList(entities []*E) []*D

	// ToEntityList 批量转换领域模型为持久化对象
	ToEntityList(domains []*D) []*E
}

// CopierMapper 基于 jinzhu/copier 的泛型 Mapper 实现
// 通过反射自动映射同名字段，支持自定义类型转换器和字段名映射
type CopierMapper[D any, E any] struct {
	converters []copier.TypeConverter
	options    copier.Option

	// 字段名映射: Domain字段名 -> Entity字段名
	// 用于处理不同 ORM 生成的字段名差异
	fieldMapping map[string]string
}

// New 创建一个新的 CopierMapper 实例
func New[D any, E any]() *CopierMapper[D, E] {
	return &CopierMapper[D, E]{
		converters:   make([]copier.TypeConverter, 0),
		fieldMapping: make(map[string]string),
		options: copier.Option{
			IgnoreEmpty: false,
			DeepCopy:    true,
		},
	}
}

// WithIgnoreEmpty 设置是否忽略空值字段
func (m *CopierMapper[D, E]) WithIgnoreEmpty(ignore bool) *CopierMapper[D, E] {
	m.options.IgnoreEmpty = ignore
	return m
}

// WithDeepCopy 设置是否深拷贝
func (m *CopierMapper[D, E]) WithDeepCopy(deep bool) *CopierMapper[D, E] {
	m.options.DeepCopy = deep
	return m
}

// WithFieldMapping 设置字段名映射
// mapping: Domain字段名 -> Entity字段名
// 例如: {"Name": "Username"} 表示 Domain.Name 对应 Entity.Username
func (m *CopierMapper[D, E]) WithFieldMapping(mapping map[string]string) *CopierMapper[D, E] {
	for k, v := range mapping {
		m.fieldMapping[k] = v
	}
	m.options.FieldNameMapping = m.buildFieldNameMapping()
	return m
}

// buildFieldNameMapping 构建 copier 的字段名映射
func (m *CopierMapper[D, E]) buildFieldNameMapping() []copier.FieldNameMapping {
	if len(m.fieldMapping) == 0 {
		return nil
	}
	mappings := make([]copier.FieldNameMapping, 0, len(m.fieldMapping)*2)
	for domainField, entityField := range m.fieldMapping {
		mappings = append(mappings, copier.FieldNameMapping{
			SrcType: new(D),
			DstType: new(E),
			Mapping: map[string]string{domainField: entityField},
		})
		mappings = append(mappings, copier.FieldNameMapping{
			SrcType: new(E),
			DstType: new(D),
			Mapping: map[string]string{entityField: domainField},
		})
	}
	return mappings
}

// RegisterConverter 注册单个类型转换器
func (m *CopierMapper[D, E]) RegisterConverter(converter copier.TypeConverter) *CopierMapper[D, E] {
	m.converters = append(m.converters, converter)
	return m
}

// RegisterConverters 批量注册类型转换器
func (m *CopierMapper[D, E]) RegisterConverters(converters []copier.TypeConverter) *CopierMapper[D, E] {
	m.converters = append(m.converters, converters...)
	return m
}

// ToDomain 将持久化对象转换为领域模型
func (m *CopierMapper[D, E]) ToDomain(entity *E) *D {
	if entity == nil {
		return nil
	}
	var domain D
	opt := m.options
	opt.Converters = m.converters
	if err := copier.CopyWithOption(&domain, entity, opt); err != nil {
		return nil
	}
	return &domain
}

// ToEntity 将领域模型转换为持久化对象
func (m *CopierMapper[D, E]) ToEntity(domain *D) *E {
	if domain == nil {
		return nil
	}
	var entity E
	opt := m.options
	opt.Converters = m.converters
	if err := copier.CopyWithOption(&entity, domain, opt); err != nil {
		return nil
	}
	return &entity
}

// ToDomainList 批量转换持久化对象为领域模型
func (m *CopierMapper[D, E]) ToDomainList(entities []*E) []*D {
	if len(entities) == 0 {
		return nil
	}
	domains := make([]*D, 0, len(entities))
	for _, entity := range entities {
		if d := m.ToDomain(entity); d != nil {
			domains = append(domains, d)
		}
	}
	return domains
}

// ToEntityList 批量转换领域模型为持久化对象
func (m *CopierMapper[D, E]) ToEntityList(domains []*D) []*E {
	if len(domains) == 0 {
		return nil
	}
	entities := make([]*E, 0, len(domains))
	for _, domain := range domains {
		if e := m.ToEntity(domain); e != nil {
			entities = append(entities, e)
		}
	}
	return entities
}
