import pika

cfg = pika.ConnectionParameters(
    host="localhost",
    credentials=pika.PlainCredentials(username="guest", password="guest"),
)

connection = pika.BlockingConnection(cfg)
channel = connection.channel()

result = channel.queue_declare(queue="", exclusive=True)
queue_name = result.method.queue
channel.queue_bind(exchange='BTCUSDT', queue=queue_name)

print(' [*] Waiting for logs. To exit press CTRL+C')
def callback(ch, method, properties, body):
    print(f" [x] {body}")

channel.basic_consume(
    queue=queue_name, on_message_callback=callback, auto_ack=True)

channel.start_consuming()