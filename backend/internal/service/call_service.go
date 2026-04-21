package service

import (
	"context"
	"errors"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type callService struct {
	callRepo        repository.CallRepository
	participantRepo repository.ParticipantRepository
	chatRepo        repository.ChatRepository
}

func NewCallService(
	callRepo repository.CallRepository,
	participantRepo repository.ParticipantRepository,
	chatRepo repository.ChatRepository,
) CallService {
	return &callService{
		callRepo:        callRepo,
		participantRepo: participantRepo,
		chatRepo:        chatRepo,
	}
}

func (s *callService) StartCall(ctx context.Context, chatID uuid.UUID, callerID uuid.UUID, calleeID uuid.UUID, callType domain.CallType) (*domain.Call, error) {
	// Проверяем, что оба участника в чате
	isCallerParticipant, err := s.participantRepo.IsParticipant(ctx, chatID, callerID)
	if err != nil || !isCallerParticipant {
		return nil, errors.New("звонящий не является участником чата")
	}

	isCalleeParticipant, err := s.participantRepo.IsParticipant(ctx, chatID, calleeID)
	if err != nil || !isCalleeParticipant {
		return nil, errors.New("вызываемый не является участником чата")
	}

	// Проверяем, нет ли уже активного звонка
	existingCall, _ := s.callRepo.GetActiveCall(ctx, chatID)
	if existingCall != nil {
		return nil, errors.New("в этом чате уже есть активный звонок")
	}

	call := &domain.Call{
		ID:        uuid.New(),
		ChatID:    chatID,
		CallerID:  callerID,
		CalleeID:  calleeID,
		Type:      callType,
		Status:    domain.CallStatusPending,
		RoomID:    uuid.New().String(),
		CreatedAt: time.Now(),
	}

	if err := s.callRepo.Create(ctx, call); err != nil {
		return nil, err
	}

	return call, nil
}

func (s *callService) AcceptCall(ctx context.Context, callID uuid.UUID) (*domain.Call, error) {
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return nil, err
	}
	if call == nil {
		return nil, errors.New("звонок не найден")
	}

	if call.Status != domain.CallStatusPending {
		return nil, errors.New("звонок уже не активен")
	}

	now := time.Now()
	call.Status = domain.CallStatusActive
	call.StartedAt = &now

	if err := s.callRepo.UpdateStatus(ctx, callID, domain.CallStatusActive); err != nil {
		return nil, err
	}

	return call, nil
}

func (s *callService) RejectCall(ctx context.Context, callID uuid.UUID) (*domain.Call, error) {
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return nil, err
	}
	if call == nil {
		return nil, errors.New("звонок не найден")
	}

	if call.Status != domain.CallStatusPending {
		return nil, errors.New("звонок уже не активен")
	}

	if err := s.callRepo.UpdateStatus(ctx, callID, domain.CallStatusRejected); err != nil {
		return nil, err
	}

	call.Status = domain.CallStatusRejected
	return call, nil
}

func (s *callService) EndCall(ctx context.Context, callID uuid.UUID) (*domain.Call, error) {
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return nil, err
	}
	if call == nil {
		return nil, errors.New("звонок не найден")
	}

	if call.Status != domain.CallStatusActive && call.Status != domain.CallStatusPending {
		return nil, errors.New("звонок уже завершен")
	}

	now := time.Now()
	if err := s.callRepo.UpdateEnded(ctx, callID, now); err != nil {
		return nil, err
	}

	call.Status = domain.CallStatusEnded
	call.EndedAt = &now
	return call, nil
}

func (s *callService) GetCallByID(ctx context.Context, callID uuid.UUID) (*domain.Call, error) {
	return s.callRepo.GetByID(ctx, callID)
}

func (s *callService) GetActiveCall(ctx context.Context, chatID uuid.UUID) (*domain.Call, error) {
	return s.callRepo.GetActiveCall(ctx, chatID)
}

func (s *callService) GetUserCalls(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Call, error) {
	return s.callRepo.GetUserCalls(ctx, userID, limit)
}
