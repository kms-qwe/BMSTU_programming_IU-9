package io.emqx.mqtt.demo;

import org.eclipse.paho.mqttv5.client.MqttClient;
import org.eclipse.paho.mqttv5.client.MqttConnectionOptions;
import org.eclipse.paho.mqttv5.common.MqttException;
import org.eclipse.paho.mqttv5.common.MqttMessage;

public class Sender {

    public static void main(String[] args) {
        String broker = "tcp://broker.emqx.io:1883";
        String clientId = "iu9_sender";
        String topic = "iu9/Krasnobaev";
        int pubQos = 1;

        try {
            MqttClient client = new MqttClient(broker, clientId);
            MqttConnectionOptions options = new MqttConnectionOptions();

            client.connect(options);

            int x = 10;
            int n = 4;
            String payload = x + "," + n;
            MqttMessage message = new MqttMessage(payload.getBytes());
            message.setQos(pubQos);
            client.publish(topic, message);

            client.disconnect();
            client.close();

        } catch (MqttException e) {
            e.printStackTrace();
        }
    }
}