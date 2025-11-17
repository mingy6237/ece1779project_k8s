package events

import (
	"inventory-manager-server/config"
	"inventory-manager-server/dto"
	"inventory-manager-server/websocket"
	"log"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
)

func InventoryUpdateHandler(msg *sarama.ConsumerMessage) {
	message := dto.Message{}
	if err := sonic.Unmarshal(msg.Value, &message); err != nil {
		log.Printf("Warning: Failed to unmarshal message: %v (continuing without Kafka consumer)", err)
		return
	}
	if message.SenderInstanceID != config.CONFIG.InstanceID {
		// This message is from another instance, broadcast it via WebSocket
		if websocket.Hub != nil {
			messageBytes, err := sonic.Marshal(message.Payload)
			if err == nil {
				websocket.Hub.Broadcast(messageBytes)
			}
		}
		return
	}
}
