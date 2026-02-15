import 'package:fetch_client/fetch_client.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:http/http.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/routes.dart';
import 'package:provider/provider.dart';

/// The Widget that configures your application.
class MyApp extends StatelessWidget {
  const MyApp({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    // Removing the line below would cause the browser URLs to not be updated when using context/route.push.
    GoRouter.optionURLReflectsImperativeAPIs = true;

    // Glue the SettingsController to the MaterialApp.
    //
    // The ListenableBuilder Widget listens to the SettingsController for changes.
    // Whenever the user updates their settings, the MaterialApp is rebuilt.
    return MultiProvider(
      providers: [
        Provider<Client>(
          create: (_) => FetchClient(
            mode: RequestMode.cors,
            credentials: RequestCredentials.cors,
          ),
          dispose: (_, client) => client.close(),
        ),
        ProxyProvider<Client, AppService>(update: (_, c, _) => AppService(c)),
      ],
      child: MaterialApp.router(
        // Providing a restorationScopeId allows the Navigator built by the
        // MaterialApp to restore the navigation stack when a user leaves and
        // returns to the app after it has been killed while running in the
        // background.
        restorationScopeId: 'app',

        supportedLocales: const [
          Locale('en', 'uk'), // English, no country code
        ],

        title: 'Admin Portal',

        // Define a light and dark color theme. Then, read the user's
        // preferred ThemeMode (light, dark, or system default) from the
        // SettingsController to display the correct theme.
        theme: ThemeData(colorSchemeSeed: Colors.amber.shade800),
        darkTheme: ThemeData.dark(),

        routerConfig: GoRouter(
          initialLocation: "/apps",
          routes: $appRoutes,
          debugLogDiagnostics: true,
        ),
      ),
    );
  }
}
