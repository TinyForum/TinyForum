'use client';

import { useState } from 'react';
import { useAuthStore } from '@/store/auth';
import SettingsSidebar from '@/components/settings/SettingsSidebar';
import DangerZone from '@/components/settings/DangerZone';
import NotificationSettings from '@/components/settings/NotificationSettings';
import ProfileSettings from '@/components/settings/ProfileSettings';
import SecuritySettings from '@/components/settings/SecuritySettings';
import AppearanceSettings from '@/components/settings/AppearanceSettings';


export type SettingsTab = 'profile' | 'security' | 'appearance' | 'notifications' | 'danger';

export default function SettingsPage() {
  const { user } = useAuthStore();
  const [activeTab, setActiveTab] = useState<SettingsTab>('profile');

  if (!user) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <div className="flex flex-col items-center gap-4">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <p className="text-base-content/60">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full min-h-screen bg-base-200/30">
      {/* 左侧菜单栏 */}
      <SettingsSidebar activeTab={activeTab} onTabChange={setActiveTab} />

      {/* 右侧内容区域 */}
      <div className="flex-1 overflow-y-auto">
        <div className="max-w-4xl mx-auto p-6 md:p-8">
          {activeTab === 'profile' && <ProfileSettings user={user} />}
          {activeTab === 'security' && <SecuritySettings user={user} />}
          {activeTab === 'appearance' && <AppearanceSettings />}
          {activeTab === 'notifications' && <NotificationSettings />}
          {activeTab === 'danger' && <DangerZone user={user} />}
        </div>
      </div>
    </div>
  );
}