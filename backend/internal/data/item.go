package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tomeai/dataflow/api/v1/flow"
	"github.com/tomeai/dataflow/internal/biz"
	"gorm.io/gorm"
)

type sourceRepo struct {
	data *Data
	log  *log.Helper
}

// QueryDataSinkById 根据ID查询单个Sink
func (s *sourceRepo) QueryDataSinkById(ctx context.Context, sinkId string) (*flow.Sink, error) {
	dataSink := new(flow.Sink)

	// 查询Sink并预加载关联的Source数据
	if err := s.data.db.Preload("Source").Where("id = ?", sinkId).First(dataSink).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warnf("Sink not found, id: %s", sinkId)
			return nil, fmt.Errorf("sink not found with id: %s", sinkId)
		}
		s.log.Errorf("query sink failed, id: %s, error: %s", sinkId, err.Error())
		return nil, err
	}

	return dataSink, nil
}

// GetAllSink 获取所有Sink数据
func (s *sourceRepo) GetAllSink(ctx context.Context) ([]*flow.Sink, error) {
	var dataSinks []*flow.Sink

	// 查询所有Sink并预加载关联的Source数据
	if err := s.data.db.Preload("Source").Find(&dataSinks).Error; err != nil {
		s.log.Errorf("query all sinks failed, error: %s", err.Error())
		return nil, err
	}

	return dataSinks, nil
}

// GetSinksBySourceId 根据Source ID获取相关的Sink列表
func (s *sourceRepo) GetSinksBySourceId(ctx context.Context, sourceId int32) ([]*flow.Sink, error) {
	var dataSinks []*flow.Sink

	if err := s.data.db.Preload("Source").Where("source_id = ?", sourceId).Find(&dataSinks).Error; err != nil {
		s.log.Errorf("query sinks by source_id failed, source_id: %d, error: %s", sourceId, err.Error())
		return nil, err
	}

	return dataSinks, nil
}

// CreateSink 创建新的Sink
func (s *sourceRepo) CreateSink(ctx context.Context, sink *flow.Sink) error {
	if err := s.data.db.Create(sink).Error; err != nil {
		s.log.Errorf("create sink failed, error: %s", err.Error())
		return err
	}

	s.log.Infof("create sink success, id: %s", sink.Id)
	return nil
}

// UpdateSink 更新Sink
func (s *sourceRepo) UpdateSink(ctx context.Context, sink *flow.Sink) error {
	if err := s.data.db.Save(sink).Error; err != nil {
		s.log.Errorf("update sink failed, id: %s, error: %s", sink.Id, err.Error())
		return err
	}

	s.log.Infof("update sink success, id: %s", sink.Id)
	return nil
}

// DeleteSink 删除Sink
func (s *sourceRepo) DeleteSink(ctx context.Context, sinkId string) error {
	if err := s.data.db.Where("id = ?", sinkId).Delete(&flow.Sink{}).Error; err != nil {
		s.log.Errorf("delete sink failed, id: %s, error: %s", sinkId, err.Error())
		return err
	}

	s.log.Infof("delete sink success, id: %s", sinkId)
	return nil
}

// QuerySourceById 根据ID查询Source
func (s *sourceRepo) QuerySourceById(ctx context.Context, sourceId int32) (*flow.Source, error) {
	source := new(flow.Source)

	if err := s.data.db.Where("id = ?", sourceId).First(source).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warnf("Source not found, id: %d", sourceId)
			return nil, fmt.Errorf("source not found with id: %d", sourceId)
		}
		s.log.Errorf("query source failed, id: %d, error: %s", sourceId, err.Error())
		return nil, err
	}

	return source, nil
}

// GetAllSource 获取所有Source数据
func (s *sourceRepo) GetAllSource(ctx context.Context) ([]*flow.Source, error) {
	var sources []*flow.Source

	if err := s.data.db.Find(&sources).Error; err != nil {
		s.log.Errorf("query all sources failed, error: %s", err.Error())
		return nil, err
	}

	return sources, nil
}

// CreateSource 创建新的Source
func (s *sourceRepo) CreateSource(ctx context.Context, source *flow.Source) error {
	if err := s.data.db.Create(source).Error; err != nil {
		s.log.Errorf("create source failed, error: %s", err.Error())
		return err
	}

	s.log.Infof("create source success, id: %d", source.Id)
	return nil
}

// UpdateSource 更新Source
func (s *sourceRepo) UpdateSource(ctx context.Context, source *flow.Source) error {
	if err := s.data.db.Save(source).Error; err != nil {
		s.log.Errorf("update source failed, id: %d, error: %s", source.Id, err.Error())
		return err
	}

	s.log.Infof("update source success, id: %d", source.Id)
	return nil
}

// DeleteSource 删除Source
func (s *sourceRepo) DeleteSource(ctx context.Context, sourceId int32) error {
	// 先检查是否有关联的Sink
	var count int64
	if err := s.data.db.Model(&flow.Sink{}).Where("source_id = ?", sourceId).Count(&count).Error; err != nil {
		s.log.Errorf("check sink count failed, source_id: %d, error: %s", sourceId, err.Error())
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete source with id %d, it has %d related sinks", sourceId, count)
	}

	if err := s.data.db.Where("id = ?", sourceId).Delete(&flow.Source{}).Error; err != nil {
		s.log.Errorf("delete source failed, id: %d, error: %s", sourceId, err.Error())
		return err
	}

	s.log.Infof("delete source success, id: %d", sourceId)
	return nil
}

// NewSourceRepo 创建新的SourceRepo实例
func NewSourceRepo(data *Data, logger log.Logger) biz.ItemRepo {
	return &sourceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
