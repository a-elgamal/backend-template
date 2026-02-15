// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'routes.dart';

// **************************************************************************
// GoRouterGenerator
// **************************************************************************

List<RouteBase> get $appRoutes => [
      $appListViewRoute,
    ];

RouteBase get $appListViewRoute => GoRouteData.$route(
      path: '/apps',
      name: 'appListView',
      factory: $AppListViewRouteExtension._fromState,
      routes: [
        GoRouteData.$route(
          path: ':appName',
          name: 'appView',
          factory: $AppViewRouteExtension._fromState,
        ),
      ],
    );

extension $AppListViewRouteExtension on AppListViewRoute {
  static AppListViewRoute _fromState(GoRouterState state) =>
      const AppListViewRoute();

  String get location => GoRouteData.$location(
        '/apps',
      );

  void go(BuildContext context) => context.go(location);

  Future<T?> push<T>(BuildContext context) => context.push<T>(location);

  void pushReplacement(BuildContext context) =>
      context.pushReplacement(location);

  void replace(BuildContext context) => context.replace(location);
}

extension $AppViewRouteExtension on AppViewRoute {
  static AppViewRoute _fromState(GoRouterState state) => AppViewRoute(
        state.pathParameters['appName']!,
        state.extra as Stored<App>?,
      );

  String get location => GoRouteData.$location(
        '/apps/${Uri.encodeComponent(appName)}',
      );

  void go(BuildContext context) => context.go(location, extra: $extra);

  Future<T?> push<T>(BuildContext context) =>
      context.push<T>(location, extra: $extra);

  void pushReplacement(BuildContext context) =>
      context.pushReplacement(location, extra: $extra);

  void replace(BuildContext context) =>
      context.replace(location, extra: $extra);
}
