package ng.hustlex.features.home.presentation

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(
    onNavigateToWallet: () -> Unit,
    onNavigateToServices: () -> Unit,
    onNavigateToRemittance: () -> Unit,
    onNavigateToProfile: () -> Unit
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Column {
                        Text(
                            "Welcome back",
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                        Text(
                            "John Doe",
                            style = MaterialTheme.typography.titleLarge,
                            fontWeight = FontWeight.Bold
                        )
                    }
                },
                actions = {
                    IconButton(onClick = { /* TODO: Notifications */ }) {
                        Icon(Icons.Default.Notifications, contentDescription = "Notifications")
                    }
                    IconButton(onClick = onNavigateToProfile) {
                        Icon(Icons.Default.Person, contentDescription = "Profile")
                    }
                }
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(24.dp)
        ) {
            // Balance Card
            item {
                WalletCard(
                    balance = "₦250,000.00",
                    onClick = onNavigateToWallet
                )
            }

            // Quick Actions
            item {
                Text(
                    "Quick Actions",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
                Spacer(modifier = Modifier.height(12.dp))
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    QuickActionItem(
                        icon = Icons.Default.Send,
                        label = "Transfer",
                        onClick = { /* TODO */ }
                    )
                    QuickActionItem(
                        icon = Icons.Default.FlightTakeoff,
                        label = "Remit",
                        onClick = onNavigateToRemittance
                    )
                    QuickActionItem(
                        icon = Icons.Default.Phone,
                        label = "Airtime",
                        onClick = { /* TODO */ }
                    )
                    QuickActionItem(
                        icon = Icons.Default.MoreHoriz,
                        label = "More",
                        onClick = { /* TODO */ }
                    )
                }
            }

            // Services Section
            item {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        "Services",
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                    TextButton(onClick = onNavigateToServices) {
                        Text("See All")
                    }
                }
            }

            item {
                LazyRow(
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    items(listOf(
                        "Cleaning" to Icons.Default.CleaningServices,
                        "Plumbing" to Icons.Default.Plumbing,
                        "Electrical" to Icons.Default.ElectricalServices,
                        "Beauty" to Icons.Default.Face,
                    )) { (name, icon) ->
                        ServiceCategoryCard(name = name, icon = icon)
                    }
                }
            }

            // Recent Transactions
            item {
                Text(
                    "Recent Transactions",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
            }

            items(listOf(
                Triple("House Cleaning", "Yesterday", "-₦15,000"),
                Triple("Airtime Purchase", "2 days ago", "-₦2,000"),
                Triple("Wallet Top-up", "3 days ago", "+₦50,000"),
            )) { (title, date, amount) ->
                TransactionItem(title = title, date = date, amount = amount)
            }
        }
    }
}

@Composable
private fun WalletCard(
    balance: String,
    onClick: () -> Unit
) {
    Card(
        onClick = onClick,
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.primary
        )
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(20.dp)
        ) {
            Text(
                "Available Balance",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onPrimary.copy(alpha = 0.8f)
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                balance,
                style = MaterialTheme.typography.headlineLarge,
                fontWeight = FontWeight.Bold,
                color = MaterialTheme.colorScheme.onPrimary
            )
            Spacer(modifier = Modifier.height(16.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                TextButton(
                    onClick = { /* TODO: Add money */ },
                    colors = ButtonDefaults.textButtonColors(
                        contentColor = MaterialTheme.colorScheme.onPrimary
                    )
                ) {
                    Icon(Icons.Default.Add, contentDescription = null)
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("Add Money")
                }
                TextButton(
                    onClick = { /* TODO: Withdraw */ },
                    colors = ButtonDefaults.textButtonColors(
                        contentColor = MaterialTheme.colorScheme.onPrimary
                    )
                ) {
                    Icon(Icons.Default.Output, contentDescription = null)
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("Withdraw")
                }
            }
        }
    }
}

@Composable
private fun QuickActionItem(
    icon: ImageVector,
    label: String,
    onClick: () -> Unit
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        FilledTonalIconButton(
            onClick = onClick,
            modifier = Modifier.size(56.dp)
        ) {
            Icon(icon, contentDescription = label)
        }
        Spacer(modifier = Modifier.height(8.dp))
        Text(
            label,
            style = MaterialTheme.typography.labelMedium
        )
    }
}

@Composable
private fun ServiceCategoryCard(
    name: String,
    icon: ImageVector
) {
    Card(
        modifier = Modifier.width(100.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Icon(
                icon,
                contentDescription = name,
                modifier = Modifier.size(32.dp),
                tint = MaterialTheme.colorScheme.primary
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                name,
                style = MaterialTheme.typography.labelMedium
            )
        }
    }
}

@Composable
private fun TransactionItem(
    title: String,
    date: String,
    amount: String
) {
    ListItem(
        headlineContent = { Text(title) },
        supportingContent = { Text(date) },
        trailingContent = {
            Text(
                amount,
                color = if (amount.startsWith("+"))
                    MaterialTheme.colorScheme.secondary
                else
                    MaterialTheme.colorScheme.error,
                fontWeight = FontWeight.Medium
            )
        }
    )
}
