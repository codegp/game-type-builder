import sys
import os
import time
import turnInformer.TurnInformer as ti
import myBot
import goless
from thrift import Thrift
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TCompactProtocol
from thrift.server import TServer
gtype = os.environ.get("GAME_TYPE")
importString = 'import ' + gtype.lower() + '.' + gtype + ' as apiclient'
exec(importString)

def NewAPIClient():
    try:
        transport = TSocket.TSocket('gamerunner', 9000)
        transport = TTransport.TBufferedTransport(transport)
        protocol = TCompactProtocol.TCompactProtocol(transport)
        client = apiclient.Client(protocol)
        transport.open()
        return client
    except Exception as e:
        time.sleep(1)
        print(e)
        return NewAPIClient()


def NewServer(addr, teamrunnerID, api, gameOverChan):
    try:
        handler = TurnInformerHandler(api, gameOverChan)
        processor = ti.Processor(handler)
        port = int(addr)
        host = "teamrunner" + teamrunnerID
        transport = TSocket.TServerSocket(host=host,port=port)
        tfactory = TTransport.TBufferedTransportFactory()
        pfactory = TCompactProtocol.TCompactProtocolFactory()
        server = TServer.TSimpleServer(processor, transport, tfactory, pfactory)
        return server
    except Exception as e:
        print(e)
        sys.exit(1)

def main(*args, **kwargs):
    addr = os.environ.get("ADDR")
    teamrunnerID = os.environ.get("TEAM_RUNNER_ID")

    apiClient = NewAPIClient()

    gameOverChan = goless.chan()
    server = NewServer(addr, teamrunnerID, apiClient, gameOverChan)
    goless.go(server.serve)
    gameOverChan.recv()
    sys.exit(0)

class TurnInformerHandler:
    def __init__(self, api, gameOverChan):
        self.api = api
        self.gameOverChan = gameOverChan

    def startTurn(self):
        print("py startTurn")
        myBot.Run(self.api)

    def gameOver(self):
        self.gameOverChan.send()
        print("py gameOver")

main(sys.argv)
