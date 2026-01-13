import express from 'express'
import {hello, test} from 'proto'
import {ChannelCredentials} from '@grpc/grpc-js'

const hc = new hello.HelloClient("worker:90", ChannelCredentials.createInsecure())
const tc = new test.TestClient("worker:90", ChannelCredentials.createInsecure())


const port = parseInt(process.env.PORT || "8080")  

const app = express()

app.get('/', (req, res) => {
    res.send('hello world')
})

app.get('/hello', async (req, res) => {
    try {
        const helloRes = await new Promise<hello.HelloRes>((resolve, reject) => {
            hc.hello(hello.HelloReq.create({msg: req.query.msg as string ?? "hi"}), (err, res) => {
                if(err) return reject(err)
                resolve(res)
            })
        })
        res.send(helloRes.msg)
    } catch(err) {
        console.error(err)
        res.status(500).end()
    }
    
})

app.listen(port, () => {
    console.log(`listening on http://localhost:${port}`);
})