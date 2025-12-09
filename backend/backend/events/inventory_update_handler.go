package events

import (
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

	// Broadcast to WebSocket for all messages (both from this instance and other instances)
	// This ensures all connected clients receive updates regardless of which node they're connected to
	if websocket.Hub != nil {
		messageBytes, err := sonic.Marshal(message.Payload)
		if err == nil {
			websocket.Hub.Broadcast(messageBytes)
			log.Printf("Broadcasted inventory update from instance %s", message.SenderInstanceID)
		} else {
			log.Printf("Failed to marshal message payload: %v", err)
		}
	}
}
