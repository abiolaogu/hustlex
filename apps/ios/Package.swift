// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "HustleX",
    platforms: [
        .iOS(.v16),
        .macOS(.v13)
    ],
    products: [
        .library(
            name: "HustleX",
            targets: ["HustleX"]
        ),
    ],
    dependencies: [
        .package(url: "https://github.com/apollographql/apollo-ios.git", from: "1.7.0"),
        .package(url: "https://github.com/pointfreeco/swift-composable-architecture", from: "1.5.0"),
    ],
    targets: [
        .target(
            name: "HustleX",
            dependencies: [
                .product(name: "Apollo", package: "apollo-ios"),
                .product(name: "ComposableArchitecture", package: "swift-composable-architecture"),
            ],
            path: "HustleX"
        ),
        .testTarget(
            name: "HustleXTests",
            dependencies: ["HustleX"]
        ),
    ]
)
