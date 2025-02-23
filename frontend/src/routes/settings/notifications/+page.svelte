<!-- src/routes/settings/notifications/+page.svelte -->
<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/config/pocketbase";
    import { Button } from "$lib/components/ui/button";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import * as Select from "$lib/components/ui/select";
    import { Separator } from "$lib/components/ui/separator";
    import { Card } from "$lib/components/ui/card";
    import { toast } from "svelte-sonner";
    import { Bell, Volume2, Mail, Clock, LayoutDashboard } from "lucide-svelte";

    let loading = true;
    let saving = false;

    let settings = {
        emailEnabled: false,
        typesEnabled: [],
        quietHoursStart: "",
        quietHoursEnd: "",
        soundEnabled: true,
        desktopNotifications: true,
        emailDigestFrequency: "daily",
        priorityLevel: "all",
    };

    const notificationTypes = [
        {
            value: "system",
            label: "System Notifications",
            description: "Important updates about your account and system",
        },
        {
            value: "task",
            label: "Task Updates",
            description: "Changes and reminders for your tasks",
        },
        {
            value: "habit",
            label: "Habit Reminders",
            description: "Daily reminders for your habits",
        },
        {
            value: "focus",
            label: "Focus Sessions",
            description: "Updates about your focus sessions",
        },
        {
            value: "calendar",
            label: "Calendar Events",
            description: "Upcoming events and meeting reminders",
        },
        {
            value: "email",
            label: "Email Notifications",
            description: "New email and message alerts",
        },
    ];

    const digestFrequencies = [
        { value: "never", label: "Never" },
        { value: "daily", label: "Daily Digest" },
        { value: "weekly", label: "Weekly Digest" },
    ];

    const priorityLevels = [
        { value: "all", label: "All Notifications" },
        { value: "important", label: "Important Only" },
        { value: "urgent", label: "Urgent Only" },
        { value: "custom", label: "Custom" },
    ];

    onMount(async () => {
        try {
            const record = await pb
                .collection("notification_settings")
                .getFirstListItem(`user = "${pb.authStore.model.id}"`);

            if (record) {
                settings = {
                    emailEnabled: record.email_enabled,
                    typesEnabled: record.types_enabled,
                    quietHoursStart: record.quiet_hours_start || "",
                    quietHoursEnd: record.quiet_hours_end || "",
                    soundEnabled: record.sound_enabled ?? true,
                    desktopNotifications: record.desktop_notifications ?? true,
                    emailDigestFrequency:
                        record.email_digest_frequency || "daily",
                    priorityLevel: record.priority_level || "all",
                };
            }
        } catch (error) {
            console.error("Error loading notification settings:", error);
            toast.error("Failed to load notification settings");
        } finally {
            loading = false;
        }
    });

    async function saveSettings() {
        saving = true;
        try {
            // Find existing record
            let record;
            try {
                record = await pb
                    .collection("notification_settings")
                    .getFirstListItem(`user = "${pb.authStore.model.id}"`);
            } catch (e) {
                // Record doesn't exist
            }

            const data = {
                user: pb.authStore.model.id,
                email_enabled: settings.emailEnabled,
                types_enabled: settings.typesEnabled,
                quiet_hours_start: settings.quietHoursStart,
                quiet_hours_end: settings.quietHoursEnd,
                sound_enabled: settings.soundEnabled,
                desktop_notifications: settings.desktopNotifications,
                email_digest_frequency: settings.emailDigestFrequency,
                priority_level: settings.priorityLevel,
            };

            if (record) {
                await pb
                    .collection("notification_settings")
                    .update(record.id, data);
            } else {
                await pb.collection("notification_settings").create(data);
            }

            toast.success("Notification settings updated successfully");
        } catch (error) {
            console.error("Error saving notification settings:", error);
            toast.error("Failed to save notification settings");
        } finally {
            saving = false;
        }
    }

    function toggleNotificationType(type: string) {
        if (settings.typesEnabled.includes(type)) {
            settings.typesEnabled = settings.typesEnabled.filter(
                (t) => t !== type,
            );
        } else {
            settings.typesEnabled = [...settings.typesEnabled, type];
        }
    }
</script>

<div class="space-y-8">
    {#if loading}
        <div class="text-center text-muted-foreground">Loading settings...</div>
    {:else}
        <!-- General Notifications -->
        <div class="space-y-4">
            <div class="flex items-center gap-2">
                <Bell class="h-5 w-5" />
                <h3 class="text-lg font-medium">General Notifications</h3>
            </div>

            <div class="grid gap-4">
                <div class="flex items-center justify-between">
                    <div class="space-y-0.5">
                        <Label>Desktop Notifications</Label>
                        <p class="text-sm text-muted-foreground">
                            Show notifications on your desktop
                        </p>
                    </div>
                    <Switch
                        checked={settings.desktopNotifications}
                        on:change={() =>
                            (settings.desktopNotifications =
                                !settings.desktopNotifications)}
                    />
                </div>

                <div class="flex items-center justify-between">
                    <div class="space-y-0.5">
                        <Label>Notification Sounds</Label>
                        <p class="text-sm text-muted-foreground">
                            Play sounds for notifications
                        </p>
                    </div>
                    <Switch
                        checked={settings.soundEnabled}
                        on:change={() =>
                            (settings.soundEnabled = !settings.soundEnabled)}
                    />
                </div>

                <div class="space-y-2">
                    <Label>Priority Level</Label>
                    <Select.Root bind:value={settings.priorityLevel}>
                        <Select.Trigger class="w-full">
                            <Select.Value placeholder="Select priority level" />
                        </Select.Trigger>
                        <Select.Content>
                            {#each priorityLevels as level}
                                <Select.Item value={level.value}
                                    >{level.label}</Select.Item
                                >
                            {/each}
                        </Select.Content>
                    </Select.Root>
                </div>
            </div>
        </div>

        <Separator />

        <!-- Email Notifications -->
        <div class="space-y-4">
            <div class="flex items-center gap-2">
                <Mail class="h-5 w-5" />
                <h3 class="text-lg font-medium">Email Notifications</h3>
            </div>

            <div class="space-y-4">
                <div class="flex items-center justify-between">
                    <div class="space-y-0.5">
                        <Label>Email Notifications</Label>
                        <p class="text-sm text-muted-foreground">
                            Receive notifications via email
                        </p>
                    </div>
                    <Switch
                        checked={settings.emailEnabled}
                        on:change={() =>
                            (settings.emailEnabled = !settings.emailEnabled)}
                    />
                </div>

                {#if settings.emailEnabled}
                    <div class="space-y-2">
                        <Label>Email Digest Frequency</Label>
                        <Select.Root bind:value={settings.emailDigestFrequency}>
                            <Select.Trigger class="w-full">
                                <Select.Value
                                    placeholder="Select digest frequency"
                                />
                            </Select.Trigger>
                            <Select.Content>
                                {#each digestFrequencies as frequency}
                                    <Select.Item value={frequency.value}
                                        >{frequency.label}</Select.Item
                                    >
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    </div>
                {/if}
            </div>
        </div>

        <Separator />

        <!-- Notification Types -->
        <div class="space-y-4">
            <div class="flex items-center gap-2">
                <LayoutDashboard class="h-5 w-5" />
                <h3 class="text-lg font-medium">Notification Types</h3>
            </div>

            <div class="grid gap-4">
                {#each notificationTypes as type}
                    <Card class="p-4">
                        <div class="flex items-center justify-between">
                            <div class="space-y-0.5">
                                <Label>{type.label}</Label>
                                <p class="text-sm text-muted-foreground">
                                    {type.description}
                                </p>
                            </div>
                            <Switch
                                checked={settings.typesEnabled.includes(
                                    type.value,
                                )}
                                on:change={() =>
                                    toggleNotificationType(type.value)}
                            />
                        </div>
                    </Card>
                {/each}
            </div>
        </div>

        <Separator />

        <!-- Quiet Hours -->
        <div class="space-y-4">
            <div class="flex items-center gap-2">
                <Clock class="h-5 w-5" />
                <h3 class="text-lg font-medium">Quiet Hours</h3>
            </div>

            <div class="space-y-4">
                <p class="text-sm text-muted-foreground">
                    Set a time range when notifications will be muted (except
                    for urgent notifications)
                </p>

                <div class="grid grid-cols-2 gap-4">
                    <div class="space-y-2">
                        <Label for="quiet-hours-start">Start Time</Label>
                        <Input
                            id="quiet-hours-start"
                            type="time"
                            bind:value={settings.quietHoursStart}
                        />
                    </div>

                    <div class="space-y-2">
                        <Label for="quiet-hours-end">End Time</Label>
                        <Input
                            id="quiet-hours-end"
                            type="time"
                            bind:value={settings.quietHoursEnd}
                        />
                    </div>
                </div>
            </div>
        </div>

        <!-- Save Button -->
        <div class="flex justify-end">
            <Button on:click={saveSettings} disabled={saving}>
                {saving ? "Saving..." : "Save Changes"}
            </Button>
        </div>
    {/if}
</div>
