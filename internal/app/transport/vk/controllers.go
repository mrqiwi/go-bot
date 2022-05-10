package vk

import (
	"context"
	"go.uber.org/zap"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type Commander interface {
	Reboot() error
	Poweroff() error
}

type VkController struct {
	logger    *zap.SugaredLogger
	bot       *api.VK
	commander Commander
}

func NewVkController(logger *zap.SugaredLogger, bot *api.VK, commander Commander) VkController {
	return VkController{
		logger:    logger,
		bot:       bot,
		commander: commander,
	}
}

func (ctrl *VkController) EventLoop() error {
	group, err := ctrl.bot.GroupsGetByID(nil)
	if err != nil {
		return err
	}

	lp, err := longpoll.NewLongPoll(ctrl.bot, group[0].ID)
	if err != nil {
		return err
	}

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		ctrl.logger.Infof("%d: %s", obj.Message.PeerID, obj.Message.Text)

		ctrl.handleMessage(strings.ToLower(obj.Message.Text), obj.Message.PeerID)
	})

	ctrl.logger.Info("Start Long Poll")

	err = lp.Run()
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *VkController) Reboot(peerID int) {
	err := ctrl.commander.Reboot()
	if err != nil {
		ctrl.logger.Errorf("Cannot run reboot: %s", err)
		return
	}

	ctrl.sendMessage(peerID, "reboot...")
}

func (ctrl *VkController) Poweroff(peerID int) {
	err := ctrl.commander.Poweroff()
	if err != nil {
		ctrl.logger.Errorf("Cannot run poweroff: %s", err)
		return
	}

	ctrl.sendMessage(peerID, "shutdowing...")
}

func (ctrl *VkController) handleMessage(text string, peerID int) {
	switch text {
	case "reset", "reboot", "restart", "перезагрузка", "рестарт":
		ctrl.Reboot(peerID)
	case "poweroff", "off", "выключение", "выкл":
		ctrl.Poweroff(peerID)
	default:
		ctrl.sendMessage(peerID, "unknown command")
	}
}

func (ctrl *VkController) sendMessage(peerID int, text string) {
	b := params.NewMessagesSendBuilder()
	b.Message(text)
	b.RandomID(0)
	b.PeerID(peerID)

	_, err := ctrl.bot.MessagesSend(b.Params)
	if err != nil {
		ctrl.logger.Errorf("Cannot send message: %s", err)
	}
}
