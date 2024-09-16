package io.emqx.mqtt.demo;

import org.eclipse.paho.mqttv5.client.IMqttToken;
import org.eclipse.paho.mqttv5.client.MqttCallback;
import org.eclipse.paho.mqttv5.client.MqttClient;
import org.eclipse.paho.mqttv5.client.MqttConnectionOptions;
import org.eclipse.paho.mqttv5.client.MqttDisconnectResponse;
import org.eclipse.paho.mqttv5.common.MqttException;
import org.eclipse.paho.mqttv5.common.MqttMessage;
import org.eclipse.paho.mqttv5.common.packet.MqttProperties;

public class Receiver {

    public static void main(String[] args) {
        String broker = "tcp://broker.emqx.io:1883";
        String clientId = "iu9_receiver";
        String topic = "iu9/Krasnobaev";
        int subQos = 1;

        try {
            MqttClient client = new MqttClient(broker, clientId);
            MqttConnectionOptions options = new MqttConnectionOptions();

            client.setCallback(new MqttCallback() {
                public void connectComplete(boolean reconnect, String serverURI) {
                    System.out.println("connected to: " + serverURI);
                }

                public void disconnected(MqttDisconnectResponse disconnectResponse) {
                    System.out.println("disconnected: " + disconnectResponse.getReasonString());
                }

                public void deliveryComplete(IMqttToken token) {
                    System.out.println("deliveryComplete: " + token.isComplete());
                }

                public void messageArrived(String topic, MqttMessage message) throws Exception {
                    System.out.println("topic: " + topic);
                    String payload = new String(message.getPayload());
                    System.out.println("message content: " + payload);

                    String[] parts = payload.split(",");
                    int x = Integer.parseInt(parts[0]);
                    int n = Integer.parseInt(parts[1]);

                    int result = (int) (n * Math.pow(x, n - 1));
                    System.out.println("Result: " + result);
                }

                public void mqttErrorOccurred(MqttException exception) {
                    System.out.println("mqttErrorOccurred: " + exception.getMessage());
                }

                public void authPacketArrived(int reasonCode, MqttProperties properties) {
                    System.out.println("authPacketArrived");
                }
            });

            client.connect(options);

            client.subscribe(topic, subQos);

        } catch (MqttException e) {
            e.printStackTrace();
        }
    }
}