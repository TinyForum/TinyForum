// components/question/RewardScoreInput.tsx
'use client';

import { CurrencyDollarIcon } from '@heroicons/react/24/outline';
import { UseFormRegister } from 'react-hook-form';
import { AskFormData } from '@/hooks/useQuestionForm';

interface RewardScoreInputProps {
  register: UseFormRegister<AskFormData>;
  rewardScore: number;
  userScore?: number;
}

export function RewardScoreInput({ register, rewardScore, userScore = 0 }: RewardScoreInputProps) {
  return (
    <div className="bg-amber-50 rounded-lg p-4 border border-amber-200">
      <label className="block text-sm font-medium text-gray-700 mb-2">
        <CurrencyDollarIcon className="w-4 h-4 inline mr-1 text-amber-600" />
        悬赏积分
      </label>
      <div className="flex items-center gap-4">
        <input
          {...register('reward_score', { 
            min: { value: 0, message: '悬赏积分不能为负数' }, 
            max: { value: 100, message: '悬赏积分不能超过100' }
          })}
          type="number"
          min="0"
          max="100"
          step="5"
          className="w-32 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
        />
        <span className="text-gray-600 text-sm">
          当前积分: <span className="font-semibold text-indigo-600">{userScore}</span>
        </span>
      </div>
      {rewardScore > 0 && (
        <p className="mt-2 text-sm text-amber-700">
          💡 悬赏 {rewardScore} 积分，回答被采纳后将扣除相应积分
        </p>
      )}
      {rewardScore === 0 && (
        <p className="mt-2 text-sm text-gray-500">
          设置悬赏积分可以吸引更多回答者
        </p>
      )}
    </div>
  );
}