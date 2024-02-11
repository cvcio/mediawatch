"""GRPC Server Module"""

import logging
from typing import Any
from collections.abc import Callable
from concurrent import futures

import grpc

from grpc_interceptor import AsyncServerInterceptor
from grpc_interceptor.exceptions import GrpcException


class ExceptionToStatusInterceptor(AsyncServerInterceptor):
    """ExceptionToStatusInterceptor class implements a new
    asyncio grpc server interceptor that converts exceptions"""

    async def intercept(
        self,
        method: Callable,
        request_or_iterator: Any,
        context: grpc.aio.ServicerContext,
        method_name: str,
    ) -> Any:
        try:
            return await method(request_or_iterator, context)
        except GrpcException as err:
            context.set_code(err.status_code)
            context.set_details(err.details)
            raise


class GRPCServer:
    """GRPC Server"""

    @property
    def instance(self) -> grpc.aio.Server:
        """return the server instance"""
        return self._server

    def __init__(
        self, host: str = "[::]", port: int = 50051, max_workers: int = 10
    ) -> None:
        options = (("grpc.so_reuseport", 1),)

        self._host = host
        self._port = port

        self._server = grpc.aio.server(
            futures.ThreadPoolExecutor(max_workers=max_workers),
            options=options,
            interceptors=[ExceptionToStatusInterceptor()],
        )

    async def stop(self) -> None:
        """
        stop the service server
        """
        logging.info("Stopping GRPC Server gracefully")
        await self._server.stop(5)

    def register_service(self, method: Callable, service: Callable, *args) -> None:
        """
        calls add_SERVICE_to_server method
        method -- the grpc method that initializes the handler
        service -- the service handler to register
        """
        method(service(args), self.instance)

    def register_service_method(self, method: Callable, service: Callable) -> None:
        """
        calls add_SERVICE_to_server method
        method -- the grpc method that initializes the handler
        service -- the service handler to register
        """
        method(service, self.instance)

    async def serve(self) -> None:
        """
        start the service server
        """
        url = f"{self._host}:{self._port}"
        logging.info("Listening at: %s", url)

        # register signals
        self._server.add_insecure_port(url)

        # start the server
        await self._server.start()
        await self._server.wait_for_termination()
