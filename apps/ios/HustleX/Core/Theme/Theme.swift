import SwiftUI

struct Theme {
    // Colors
    struct Colors {
        static let primary = Color("Primary", bundle: nil)
        static let secondary = Color("Secondary", bundle: nil)
        static let accent = Color("Accent", bundle: nil)
        static let background = Color("Background", bundle: nil)
        static let surface = Color("Surface", bundle: nil)
        static let error = Color.red
        static let success = Color.green

        // Fallback colors if assets not available
        static let primaryFallback = Color(hex: "4F46E5")
        static let secondaryFallback = Color(hex: "10B981")
        static let accentFallback = Color(hex: "F59E0B")
    }

    // Typography
    struct Typography {
        static let displayLarge = Font.system(size: 57, weight: .bold)
        static let displayMedium = Font.system(size: 45, weight: .bold)
        static let displaySmall = Font.system(size: 36, weight: .bold)
        static let headlineLarge = Font.system(size: 32, weight: .semibold)
        static let headlineMedium = Font.system(size: 28, weight: .semibold)
        static let headlineSmall = Font.system(size: 24, weight: .semibold)
        static let titleLarge = Font.system(size: 22, weight: .medium)
        static let titleMedium = Font.system(size: 16, weight: .medium)
        static let titleSmall = Font.system(size: 14, weight: .medium)
        static let bodyLarge = Font.system(size: 16, weight: .regular)
        static let bodyMedium = Font.system(size: 14, weight: .regular)
        static let bodySmall = Font.system(size: 12, weight: .regular)
        static let labelLarge = Font.system(size: 14, weight: .medium)
        static let labelMedium = Font.system(size: 12, weight: .medium)
        static let labelSmall = Font.system(size: 11, weight: .medium)
    }

    // Spacing
    struct Spacing {
        static let xs: CGFloat = 4
        static let sm: CGFloat = 8
        static let md: CGFloat = 16
        static let lg: CGFloat = 24
        static let xl: CGFloat = 32
        static let xxl: CGFloat = 48
    }

    // Corner Radius
    struct CornerRadius {
        static let sm: CGFloat = 4
        static let md: CGFloat = 8
        static let lg: CGFloat = 12
        static let xl: CGFloat = 16
        static let full: CGFloat = 9999
    }
}

extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3:
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6:
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8:
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (255, 0, 0, 0)
        }
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
