// 'use client';

// import { useState, useEffect } from 'react';
// import { useForm } from 'react-hook-form';
// import { z } from 'zod';
// import { zodResolver } from '@hookform/resolvers/zod';
// import { authApi } from '@/lib/api';
// import { Mail, Lock, AlertTriangle, RefreshCw, Trash2, CheckCircle } from 'lucide-react';
// import { useTranslations } from 'next-intl';

// const recoverySchema = z.object({
//   email: z.string().email('请输入有效的邮箱'),
//   password: z.string().min(1, '请输入密码'),
// });

// type RecoveryForm = z.infer<typeof recoverySchema>;

// interface DeletionStatus {
//   is_deleted: boolean;
//   deleted_at?: string;
//   can_restore: boolean;
//   remaining_days?: number;
// }

// export default function AccountRecovery() {
//   const t = useTranslations('Auth');
//   const [step, setStep] = useState<'form' | 'status'>('form');
//   const [isLoading, setIsLoading] = useState(false);
//   const [error, setError] = useState<string | null>(null);
//   const [deletionStatus, setDeletionStatus] = useState<DeletionStatus | null>(null);
//   const [userEmail, setUserEmail] = useState('');

//   const {
//     register,
//     handleSubmit,
//     formState: { errors },
//   } = useForm<RecoveryForm>({
//     resolver: zodResolver(recoverySchema),
//   });

//   // 检查账户删除状态
//   const checkDeletionStatus = async (email: string, password: string) => {
//     setIsLoading(true);
//     setError(null);

//     try {
//       // 先登录获取 token
//       const loginResponse = await authApi.login({ email, password });

//       if (loginResponse.data.code === 200) {
//         // 获取删除状态
//         const statusResponse = await authApi.getDeletionStatus();
//         setDeletionStatus(statusResponse.data.data);
//         setUserEmail(email);
//         setStep('status');
//       }
//     } catch (err: any) {
//       setError(err?.response?.data?.message || '验证失败，请检查邮箱和密码');
//     } finally {
//       setIsLoading(false);
//     }
//   };

//   // 恢复账户
//   const handleRestore = async () => {
//     setIsLoading(true);
//     setError(null);

//     try {
//       await authApi.cancelDeletion();
//       // 刷新状态
//       const statusResponse = await authApi.getDeletionStatus();
//       setDeletionStatus(statusResponse.data.data);
//       alert('账户已成功恢复！');
//     } catch (err: any) {
//       setError(err?.response?.data?.message || '恢复失败，请重试');
//     } finally {
//       setIsLoading(false);
//     }
//   };

//   // 永久删除
//   const handlePermanentDelete = async () => {
//     const confirm = window.confirm('⚠️ 警告：此操作将永久删除您的账户，所有数据无法恢复！\n\n确认永久删除吗？');
//     if (!confirm) return;

//     setIsLoading(true);
//     setError(null);

//     try {
//       await authApi.confirmDeletion({ confirm: 'PERMANENT_DELETE' });
//       alert('账户已永久删除');
//       // 登出并返回首页
//       await authApi.logout();
//       window.location.href = '/';
//     } catch (err: any) {
//       setError(err?.response?.data?.message || '删除失败，请重试');
//     } finally {
//       setIsLoading(false);
//     }
//   };

//   const onSubmit = (data: RecoveryForm) => {
//     checkDeletionStatus(data.email, data.password);
//   };

//   const resetForm = () => {
//     setStep('form');
//     setDeletionStatus(null);
//     setError(null);
//     setUserEmail('');
//   };

//   return (
//     <div className="mt-8 pt-6 border-t border-base-300">
//       <div className="text-center mb-4">
//         <h3 className="text-lg font-semibold flex items-center justify-center gap-2">
//           <AlertTriangle className="w-5 h-5 text-warning" />
//           {t('account_management')}
//         </h3>
//         <p className="text-sm text-base-content/60 mt-1">
//           {t('check_or_recover_account')}
//         </p>
//       </div>

//       {step === 'form' ? (
//         <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
//           <div className="form-control">
//             <label className="label pb-1">
//               <span className="label-text font-medium">{t('email')}</span>
//             </label>
//             <div className="relative">
//               <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
//               <input
//                 {...register('email')}
//                 type="email"
//                 placeholder="your@email.com"
//                 className={`input input-bordered w-full pl-10 ${errors.email ? 'input-error' : ''}`}
//                 autoComplete="email"
//               />
//             </div>
//             {errors.email && (
//               <label className="label pt-1">
//                 <span className="label-text-alt text-error">{errors.email.message}</span>
//               </label>
//             )}
//           </div>

//           <div className="form-control">
//             <label className="label pb-1">
//               <span className="label-text font-medium">{t('password')}</span>
//             </label>
//             <div className="relative">
//               <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
//               <input
//                 {...register('password')}
//                 type="password"
//                 placeholder="••••••••"
//                 className={`input input-bordered w-full pl-10 ${errors.password ? 'input-error' : ''}`}
//                 autoComplete="current-password"
//               />
//             </div>
//             {errors.password && (
//               <label className="label pt-1">
//                 <span className="label-text-alt text-error">{errors.password.message}</span>
//               </label>
//             )}
//           </div>

//           {error && (
//             <div className="alert alert-error py-2">
//               <span className="text-sm">{error}</span>
//             </div>
//           )}

//           <button
//             type="submit"
//             className="btn btn-outline btn-warning w-full gap-2"
//             disabled={isLoading}
//           >
//             {isLoading ? (
//               <span className="loading loading-spinner loading-sm" />
//             ) : (
//               <RefreshCw className="w-4 h-4" />
//             )}
//             {t('check_status')}
//           </button>
//         </form>
//       ) : (
//         <div className="space-y-4">
//           {deletionStatus?.is_deleted ? (
//             <div className="alert alert-warning">
//               <AlertTriangle className="w-5 h-5" />
//               <div className="flex-1">
//                 <p className="font-semibold">{t('account_marked_for_deletion')}</p>
//                 <p className="text-sm">
//                   {deletionStatus.deleted_at && (
//                     <>
//                       删除时间：{new Date(deletionStatus.deleted_at).toLocaleDateString()}
//                       <br />
//                     </>
//                   )}
//                   {deletionStatus.can_restore && (
//                     <span className="text-error">
//                       剩余 {deletionStatus.remaining_days} 天可恢复，逾期将永久删除
//                     </span>
//                   )}
//                 </p>
//               </div>
//             </div>
//           ) : (
//             <div className="alert alert-success">
//               <CheckCircle className="w-5 h-5" />
//               <div>
//                 <p className="font-semibold">{t('account_normal')}</p>
//                 <p className="text-sm">{userEmail}</p>
//               </div>
//             </div>
//           )}

//           {deletionStatus?.is_deleted && deletionStatus.can_restore && (
//             <div className="flex gap-3">
//               <button
//                 onClick={handleRestore}
//                 className="btn btn-success flex-1 gap-2"
//                 disabled={isLoading}
//               >
//                 <RefreshCw className="w-4 h-4" />
//                 {t('restore_account')}
//               </button>
//               <button
//                 onClick={handlePermanentDelete}
//                 className="btn btn-error flex-1 gap-2"
//                 disabled={isLoading}
//               >
//                 <Trash2 className="w-4 h-4" />
//                 {t('permanent_delete')}
//               </button>
//             </div>
//           )}

//           <button
//             onClick={resetForm}
//             className="btn btn-ghost btn-sm w-full"
//           >
//             {t('check_another_account')}
//           </button>
//         </div>
//       )}
//     </div>
//   );
// }
