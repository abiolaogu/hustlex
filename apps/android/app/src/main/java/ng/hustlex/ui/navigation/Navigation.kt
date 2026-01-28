package ng.hustlex.ui.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import ng.hustlex.features.auth.presentation.LoginScreen
import ng.hustlex.features.home.presentation.HomeScreen
import ng.hustlex.features.wallet.presentation.WalletScreen
import ng.hustlex.features.services.presentation.ServicesScreen
import ng.hustlex.features.remittance.presentation.RemittanceScreen
import ng.hustlex.features.profile.presentation.ProfileScreen

sealed class Screen(val route: String) {
    object Login : Screen("login")
    object Home : Screen("home")
    object Wallet : Screen("wallet")
    object Services : Screen("services")
    object Remittance : Screen("remittance")
    object Profile : Screen("profile")
}

@Composable
fun HustleXNavHost(
    navController: NavHostController = rememberNavController(),
    startDestination: String = Screen.Login.route
) {
    NavHost(
        navController = navController,
        startDestination = startDestination
    ) {
        composable(Screen.Login.route) {
            LoginScreen(
                onLoginSuccess = {
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.Login.route) { inclusive = true }
                    }
                }
            )
        }

        composable(Screen.Home.route) {
            HomeScreen(
                onNavigateToWallet = { navController.navigate(Screen.Wallet.route) },
                onNavigateToServices = { navController.navigate(Screen.Services.route) },
                onNavigateToRemittance = { navController.navigate(Screen.Remittance.route) },
                onNavigateToProfile = { navController.navigate(Screen.Profile.route) }
            )
        }

        composable(Screen.Wallet.route) {
            WalletScreen(
                onNavigateBack = { navController.popBackStack() }
            )
        }

        composable(Screen.Services.route) {
            ServicesScreen(
                onNavigateBack = { navController.popBackStack() }
            )
        }

        composable(Screen.Remittance.route) {
            RemittanceScreen(
                onNavigateBack = { navController.popBackStack() }
            )
        }

        composable(Screen.Profile.route) {
            ProfileScreen(
                onNavigateBack = { navController.popBackStack() },
                onLogout = {
                    navController.navigate(Screen.Login.route) {
                        popUpTo(0) { inclusive = true }
                    }
                }
            )
        }
    }
}
