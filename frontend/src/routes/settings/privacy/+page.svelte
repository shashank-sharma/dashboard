<!-- src/routes/settings/privacy/+page.svelte -->
<script lang="ts">
    import { pb } from "$lib/config/pocketbase";
    import { Switch } from "$lib/components/ui/switch";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { toast } from "svelte-sonner";
    import * as Select from "$lib/components/ui/select";
    import { Shield, Lock, Eye, Key } from "lucide-svelte";

    // Example privacy settings
    let settings = {
        twoFactorEnabled: false,
        activityLogging: true,
        dataSharing: "minimal",
        loginNotifications: true,
        publicProfile: false,
        passwordLastChanged: new Date("2024-01-01"),
        securityLevel: "high",
    };

    const securityLevels = [
        { value: "low", label: "Basic" },
        { value: "medium", label: "Standard" },
        { value: "high", label: "Enhanced" },
        { value: "custom", label: "Custom" },
    ];

    async function saveSettings() {
        try {
            // Example of saving to PocketBase
            // await pb.collection('user_privacy_settings').update(id, settings);
            toast.success("Privacy settings updated successfully");
        } catch (error) {
            toast.error("Failed to update privacy settings");
        }
    }

    function resetSecuritySettings() {
        // Implementation for resetting security settings
        toast.success("Security settings have been reset");
    }
</script>

<div class="space-y-8">
    <!-- Security Settings -->
    <div class="space-y-4">
        <div class="flex items-center gap-2">
            <Shield class="h-5 w-5" />
            <h3 class="text-lg font-medium">Security Settings</h3>
        </div>

        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                    <Label>Two-Factor Authentication</Label>
                    <p class="text-sm text-muted-foreground">
                        Add an extra layer of security to your account
                    </p>
                </div>
                <Switch
                    checked={settings.twoFactorEnabled}
                    on:change={() =>
                        (settings.twoFactorEnabled =
                            !settings.twoFactorEnabled)}
                />
            </div>

            <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                    <Label>Login Notifications</Label>
                    <p class="text-sm text-muted-foreground">
                        Get notified of new login attempts
                    </p>
                </div>
                <Switch
                    checked={settings.loginNotifications}
                    on:change={() =>
                        (settings.loginNotifications =
                            !settings.loginNotifications)}
                />
            </div>

            <div class="space-y-2">
                <Label>Security Level</Label>
                <Select.Root value={settings.securityLevel}>
                    <Select.Trigger class="w-full">
                        <Select.Value placeholder="Select security level" />
                    </Select.Trigger>
                    <Select.Content>
                        {#each securityLevels as level}
                            <Select.Item value={level.value}>
                                {level.label}
                            </Select.Item>
                        {/each}
                    </Select.Content>
                </Select.Root>
            </div>
        </div>
    </div>

    <!-- Privacy Settings -->
    <div class="space-y-4">
        <div class="flex items-center gap-2">
            <Eye class="h-5 w-5" />
            <h3 class="text-lg font-medium">Privacy Settings</h3>
        </div>

        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                    <Label>Activity Logging</Label>
                    <p class="text-sm text-muted-foreground">
                        Track your account activity for security
                    </p>
                </div>
                <Switch
                    checked={settings.activityLogging}
                    on:change={() =>
                        (settings.activityLogging = !settings.activityLogging)}
                />
            </div>

            <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                    <Label>Public Profile</Label>
                    <p class="text-sm text-muted-foreground">
                        Make your profile visible to others
                    </p>
                </div>
                <Switch
                    checked={settings.publicProfile}
                    on:change={() =>
                        (settings.publicProfile = !settings.publicProfile)}
                />
            </div>

            <div class="space-y-2">
                <Label>Data Sharing</Label>
                <Select.Root value={settings.dataSharing}>
                    <Select.Trigger class="w-full">
                        <Select.Value placeholder="Select data sharing level" />
                    </Select.Trigger>
                    <Select.Content>
                        <Select.Item value="none">No sharing</Select.Item>
                        <Select.Item value="minimal">Minimal</Select.Item>
                        <Select.Item value="standard">Standard</Select.Item>
                        <Select.Item value="full">Full</Select.Item>
                    </Select.Content>
                </Select.Root>
            </div>
        </div>
    </div>

    <!-- Password Management -->
    <div class="space-y-4">
        <div class="flex items-center gap-2">
            <Key class="h-5 w-5" />
            <h3 class="text-lg font-medium">Password Management</h3>
        </div>

        <div class="space-y-4">
            <div class="space-y-2">
                <Label>Last Password Change</Label>
                <p class="text-sm text-muted-foreground">
                    {settings.passwordLastChanged.toLocaleDateString()}
                </p>
            </div>

            <div class="flex gap-4">
                <Button variant="outline" on:click={resetSecuritySettings}>
                    Reset Security Settings
                </Button>
                <Button on:click={saveSettings}>Save Changes</Button>
            </div>
        </div>
    </div>
</div>
