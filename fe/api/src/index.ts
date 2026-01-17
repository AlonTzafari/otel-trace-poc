import { startOtel } from './tracing'
startOtel()


import express from 'express'
import { hello, test } from 'proto'
import { ChannelCredentials } from '@grpc/grpc-js'
import { Kafka } from 'kafkajs'


const workerHost = process.env.WORKER_HOST || "worker"
const workerPort = process.env.WORKER_PORT || "90"
const workerAddr = `${workerHost}:${workerPort}`

const hc = new hello.HelloClient(workerAddr, ChannelCredentials.createInsecure())
const tc = new test.TestClient(workerAddr, ChannelCredentials.createInsecure())

const kafka = new Kafka({
    brokers: ["kafka:9092"],
})

const consumer = kafka.consumer({ groupId: "api", allowAutoTopicCreation: true })
consumer.connect()


const port = parseInt(process.env.PORT || "8080")

const app = express()

app.get('/', (req, res) => {
    res.send('hello world')
})

app.get('/hello', async (req, res) => {
    try {
        const helloRes = await new Promise<hello.HelloRes>((resolve, reject) => {
            hc.hello(hello.HelloReq.create({ msg: req.query.msg as string ?? "hi" }), (err, res) => {
                if (err) return reject(err)
                resolve(res)
            })
        })
        res.send(helloRes.msg)
    } catch (err) {
        console.error(err)
        res.status(500).end()
    }

})

async function start() {
    console.log("consumer.connect");
    await consumer.connect()
    console.log("consumer.subscribe");
    await consumer.subscribe({
        topic: "diceroll",
        fromBeginning: true,
    })
    console.log("consumer.run");
    consumer.run({
        eachMessage: async ({topic, partition, message}) => {
            if (!message.value) {
                throw new Error("invalid message")
            }
            const event = hello.DiceRollEvent.decode(message.value)
            console.log("dice roll received", event)
        }
    })
    
    console.log("app.listen");
    const server = app.listen(port, () => {
        console.log(`listening on http://localhost:${port}`);
    })

    server.on("error", (err) => {
        console.error(err)
    })

    server.on("close", (e) => {
        console.info(e)
    })
}

start()

