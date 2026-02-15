import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:go_router/go_router.dart';
import 'package:portal/constants.dart';
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/routes.dart';
import 'package:portal/src/apps/new_app_view.dart';
import 'package:portal/src/stored/stored.dart';
import 'package:portal/src/widgets/error_dialog.dart';

class AppListView extends StatefulWidget {
  final AppService _appService;

  const AppListView(this._appService, {super.key});

  @override
  State<AppListView> createState() => _AppListViewState();
}

class _AppListViewState extends State<AppListView> {
  List<Stored<App>> _apps = [];
  var loading = true;
  var showDeleted = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    loadApps();
  }

  void loadApps() async {
    try {
      var fetchedApps =
          await widget._appService.listApps(showDeleted: showDeleted);
      setState(() => _apps = fetchedApps);
    } catch (error, stackTrace) {
      if (mounted) {
        SchedulerBinding.instance.addPostFrameCallback((_) => showError(
            context,
            error,
            stackTrace,
            "An error occured while fetching apps list from the backend"));
      }
    } finally {
      setState(() => loading = false);
    }
  }

  void addApp() async {
    final name = await showModalBottomSheet(
        context: context,
        builder: (ctx) => Padding(
              padding: const EdgeInsets.all(8.0),
              child: NewAppView(
                onSave: (name) => ctx.pop(name),
              ),
            ));

    if (name == null) {
      return;
    }
    setState(() {
      loading = true;
    });

    try {
      var p = await widget._appService.addApp(name);
      setState(() {
        _apps.add(p);
      });
      showAppDetail(_apps.length - 1);
    } catch (error, stackTrace) {
      if (mounted) {
        SchedulerBinding.instance.addPostFrameCallback((_) => showError(context,
            error, stackTrace, "An error occured while adding app: '$name'"));
      }
    } finally {
      setState(() {
        loading = false;
      });
    }
  }

  void showAppDetail(int i) async {
    final Stored<App>? app =
        await AppViewRoute(_apps[i].id, _apps[i]).push(context);
    if (app != null && app.content.disabled && !showDeleted) {
      _apps.removeAt(i);
    } else if (app != null) {
      _apps.replaceRange(i, i + 1, [app]);
    }
  }

  void toggleDeletedApps() {
    setState(() {
      loading = true;
      showDeleted = !showDeleted;
    });
    loadApps();
  }

  @override
  Widget build(BuildContext context) {
    Widget body = Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        const Center(child: Text("No Apps to show!")),
        const SizedBox(
          height: 15,
        ),
        if (!showDeleted)
          ElevatedButton.icon(
            onPressed: toggleDeletedApps,
            icon: const Icon(Icons.visibility),
            label: const Text("Show deleted"),
          )
      ],
    );

    final theme = Theme.of(context);

    if (loading) {
      body = const Center(child: CircularProgressIndicator());
    } else if (_apps.isNotEmpty) {
      body = Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: maxDisplayWidth),
          child: Column(
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  IconButton(
                    onPressed: toggleDeletedApps,
                    icon: Icon(
                        showDeleted ? Icons.visibility_off : Icons.visibility),
                  )
                ],
              ),
              Expanded(
                child: ListView.builder(
                  itemCount: _apps.length,
                  itemBuilder: (ctx, i) => ListTile(
                    leading: CircleAvatar(
                      child: Text(
                        "${i + 1}",
                        style: theme.textTheme.labelLarge,
                      ),
                    ),
                    title: Text(
                      _apps[i].id,
                      style: _apps[i].content.disabled
                          ? theme.textTheme.labelLarge!.copyWith(
                              color: theme.colorScheme.error,
                              decoration: TextDecoration.lineThrough)
                          : theme.textTheme.labelLarge,
                    ),
                    trailing: const Icon(Icons.arrow_forward),
                    onTap: () => showAppDetail(i),
                  ),
                ),
              ),
            ],
          ),
        ),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text("Apps"),
            IconButton(
              onPressed: () {
                setState(() => loading = true);
                loadApps();
              },
              icon: const Icon(Icons.refresh),
            )
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.large(
        onPressed: addApp,
        child: const Icon(Icons.add),
      ),
      body: body,
    );
  }
}
