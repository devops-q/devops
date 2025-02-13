FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS base
USER $APP_UID
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["itu-minitwit-cs.csproj", "./"]
RUN dotnet restore "itu-minitwit-cs.csproj"
COPY . .
WORKDIR "/src/"
RUN dotnet build "itu-minitwit-cs.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "itu-minitwit-cs.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app

COPY --from=publish /app/publish .
COPY templates /app/templates

VOLUME ["/app/data"]
ENTRYPOINT ["dotnet", "itu-minitwit-cs.dll"]