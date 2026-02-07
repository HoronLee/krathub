package mapper

import (
	"github.com/jinzhu/copier"
)

// ProtoMapper 专门用于 Protobuf 消息与领域模型之间的转换
// P = Protobuf message (API层)
// D = Domain model (biz层)
type ProtoMapper[P any, D any] struct {
	converters []copier.TypeConverter
	options    copier.Option
}

// NewProtoMapper 创建一个新的 ProtoMapper 实例
// 默认注册 timestamppb 转换器
func NewProtoMapper[P any, D any]() *ProtoMapper[P, D] {
	m := &ProtoMapper[P, D]{
		converters: make([]copier.TypeConverter, 0),
		options: copier.Option{
			IgnoreEmpty: false,
			DeepCopy:    true,
		},
	}
	// 默认注册 protobuf 常用转换器
	m.converters = append(m.converters, NewTimestamppbConverterPair()...)
	m.converters = append(m.converters, NewStringPointerConverterPair()...)
	m.converters = append(m.converters, NewInt64PointerConverterPair()...)
	return m
}

// RegisterConverter 注册单个类型转换器
func (m *ProtoMapper[P, D]) RegisterConverter(converter copier.TypeConverter) *ProtoMapper[P, D] {
	m.converters = append(m.converters, converter)
	return m
}

// RegisterConverters 批量注册类型转换器
func (m *ProtoMapper[P, D]) RegisterConverters(converters []copier.TypeConverter) *ProtoMapper[P, D] {
	m.converters = append(m.converters, converters...)
	return m
}

// ToDomain 将 Protobuf 消息转换为领域模型
func (m *ProtoMapper[P, D]) ToDomain(proto *P) *D {
	if proto == nil {
		return nil
	}
	var domain D
	opt := m.options
	opt.Converters = m.converters
	if err := copier.CopyWithOption(&domain, proto, opt); err != nil {
		return nil
	}
	return &domain
}

// ToProto 将领域模型转换为 Protobuf 消息
func (m *ProtoMapper[P, D]) ToProto(domain *D) *P {
	if domain == nil {
		return nil
	}
	var proto P
	opt := m.options
	opt.Converters = m.converters
	if err := copier.CopyWithOption(&proto, domain, opt); err != nil {
		return nil
	}
	return &proto
}

// ToDomainList 批量转换 Protobuf 消息为领域模型
func (m *ProtoMapper[P, D]) ToDomainList(protos []*P) []*D {
	if len(protos) == 0 {
		return nil
	}
	domains := make([]*D, 0, len(protos))
	for _, proto := range protos {
		if d := m.ToDomain(proto); d != nil {
			domains = append(domains, d)
		}
	}
	return domains
}

// ToProtoList 批量转换领域模型为 Protobuf 消息
func (m *ProtoMapper[P, D]) ToProtoList(domains []*D) []*P {
	if len(domains) == 0 {
		return nil
	}
	protos := make([]*P, 0, len(domains))
	for _, domain := range domains {
		if p := m.ToProto(domain); p != nil {
			protos = append(protos, p)
		}
	}
	return protos
}
