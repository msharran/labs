import asyncio

import click
import websockets

clients = set()


async def handler(ws):
    clients.add(ws)
    try:
        async for msg in ws:
            for client in clients:
                if client != ws:
                    await client.send(msg)
    finally:
        clients.remove(ws)


async def main():
    async with websockets.serve(handler, "0.0.0.0", 8765):
        await asyncio.Future()


@click.command()
def root():
    asyncio.run(main())
