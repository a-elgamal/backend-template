import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:portal/constants.dart';
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/stored/stored.dart';
import 'package:portal/src/widgets/confirmation_dialog.dart';
import 'package:portal/src/widgets/error_dialog.dart';

class AppView extends StatefulWidget {
  final AppService _appService;
  final String appName;
  final Stored<App>? app;

  const AppView(this._appService, this.appName, {super.key, this.app});

  @override
  AppViewState createState() => AppViewState();
}

class AppViewState extends State<AppView> {
  Stored<App>? app;
  bool loading = false;

  @override
  void initState() {
    super.initState();
    if (widget.app != null) {
      app = widget.app!;
    } else {
      loading = true;
    }
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (app == null) {
      loadApp();
    }
  }

  void loadApp() async {
    try {
      final result = await widget._appService.getApp(widget.appName);
      setState(() {
        app = result;
        loading = false;
      });
    } catch (error, stackTrace) {
      if (mounted) {
        SchedulerBinding.instance.addPostFrameCallback((_) => showError(context,
            error, stackTrace, "Failed to fetch app: '${widget.appName}"));
      }
    }
  }

  void toggleDeleted() async {
    final yes = await showConfirmationDialog(
      context,
      title: app!.content.disabled
          ? "Restore ${widget.appName}?"
          : "Delete ${widget.appName}?",
      text: app!.content.disabled
          ? "All requests from this app can resume with the same API Key"
          : "All requests from this app will stop working immediately.",
    );
    if (yes == null || !yes) {
      return;
    }

    setState(() {
      loading = true;
    });
    try {
      await widget._appService
          .setDisabled(widget.appName, !app!.content.disabled);
      setState(() {
        app = Stored.copy(app!,
            content: App.copy(app!.content, disabled: !app!.content.disabled));
      });
    } catch (error, stackTrace) {
      if (mounted) {
        SchedulerBinding.instance.addPostFrameCallback((_) => showError(
            context,
            error,
            stackTrace,
            "An error occured while deleting app: '${widget.appName}'"));
      }
    } finally {
      setState(() {
        loading = false;
      });
    }
  }

  void resetAPIKey() async {
    final yes = await showConfirmationDialog(
      context,
      title: "Generate New API Key for ${widget.appName}?",
      text:
          "The current API Key will be invalidated and any ongoing requests using it will fail!",
    );
    if (yes == null || !yes) {
      return;
    }

    setState(() {
      loading = true;
    });
    try {
      await widget._appService.resetAPIKey(widget.appName);
      loadApp();
    } catch (error, stackTrace) {
      setState(() {
        loading = false;
      });
      if (mounted) {
        SchedulerBinding.instance.addPostFrameCallback((_) => showError(
            context,
            error,
            stackTrace,
            "An error occured while deleting app: '${widget.appName}'"));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    Widget body = const Center(child: CircularProgressIndicator());

    const rowHeight = 40.0;

    final theme = Theme.of(context);

    if (!loading && app != null) {
      body = Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Padding(
            padding: const EdgeInsets.all(5.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                SizedBox(
                  height: rowHeight,
                  child: Center(
                    child: Text(
                      "API Key:",
                      style: theme.textTheme.titleMedium,
                    ),
                  ),
                ),
                SizedBox(
                  height: rowHeight,
                  child: Center(
                    child: Text(
                      "Created By:",
                      style: theme.textTheme.titleMedium,
                    ),
                  ),
                ),
                SizedBox(
                  height: rowHeight,
                  child: Center(
                    child: Text(
                      "Created At:",
                      style: theme.textTheme.titleMedium,
                    ),
                  ),
                ),
                SizedBox(
                  height: rowHeight,
                  child: Center(
                    child: Text(
                      "Last Modified By:",
                      style: theme.textTheme.titleMedium,
                    ),
                  ),
                ),
                SizedBox(
                  height: rowHeight,
                  child: Center(
                    child: Text(
                      "Last Modified At:",
                      style: theme.textTheme.titleMedium,
                    ),
                  ),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(5.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(
                  height: rowHeight,
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      SelectableText(
                        app!.content.apiKey,
                        key: const Key("api-key"),
                      ),
                      if (!app!.content.disabled)
                        IconButton(
                            onPressed: resetAPIKey,
                            icon: const Icon(Icons.create)),
                    ],
                  ),
                ),
                SizedBox(
                  height: rowHeight,
                  child: Center(child: Text(app!.createdBy)),
                ),
                SizedBox(
                  height: rowHeight,
                  child:
                      Center(child: Text(app!.createdAt.toLocal().toString())),
                ),
                SizedBox(
                  height: rowHeight,
                  child: Center(child: Text(app!.modifiedBy)),
                ),
                SizedBox(
                  height: rowHeight,
                  child:
                      Center(child: Text(app!.modifiedAt.toLocal().toString())),
                ),
              ],
            ),
          ),
        ],
      );
    }

    return Scaffold(
      appBar: AppBar(
        leading: BackButton(
          onPressed: () => Navigator.of(context).pop(app),
        ),
        title: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SelectableText(
              widget.appName,
              style: app?.content.disabled ?? false
                  ? theme.textTheme.titleLarge!.copyWith(
                      color: theme.colorScheme.error,
                      decoration: TextDecoration.lineThrough)
                  : theme.textTheme.titleLarge,
            ),
            if (!(app?.content.disabled ?? false))
              IconButton(
                onPressed: () {
                  setState(() => loading = true);
                  loadApp();
                },
                icon: const Icon(Icons.refresh),
              )
          ],
        ),
        actions: [
          IconButton(
            color: app?.content.disabled ?? false
                ? null
                : Theme.of(context).colorScheme.error,
            onPressed: toggleDeleted,
            icon: app?.content.disabled ?? false
                ? const Icon(Icons.restore)
                : const Icon(Icons.delete),
          ),
        ],
      ),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: maxDisplayWidth),
          child: body,
        ),
      ),
    );
  }
}
