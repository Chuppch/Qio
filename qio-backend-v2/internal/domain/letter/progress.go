package letter

import "time"

// 投递进度计算。
//
// 迁移自 v1 的 ProgressUtils.getProgress。算法保持一致，修正了两处问题：
//
//  1. v1 用 `progress < 0 || progress > 10000 ? 10000 : progress` 做钳制，
//     负值会被错误地当成已送达。这里负值钳到 0。
//  2. v1 直接读 letter.reduceTime / speedRate 两个字符串字段做解析，解析失败会 panic。
//     强类型化后由仓储层负责转换，本层不再处理字符串。
//
// 另外 v1 在计算进度的同时顺手改写了 DeliveryTime，属于把两件事塞进一个函数。
// 这里拆成 Progress 与 ActualDeliveryAt 两个纯函数，调用方按需组合。

// Progress 返回信件在 now 时刻的投递进度，万分制，取值范围 [0, ProgressScale]。
//
// 计算方式：以「未加速的总时长」为分母，把道具减免掉的时长折算成已完成的进度，
// 与真实流逝时间相加后得到当前进度。
func (l *Letter) Progress(now time.Time) int64 {
	originalTotal := l.ExpectedDeliveryAt.Sub(l.CreatedAt)
	if originalTotal <= 0 {
		return ProgressScale
	}

	newTotal := l.acceleratedTotal(originalTotal)
	// 道具减免掉的时长，直接折算为已完成部分
	saved := originalTotal - newTotal
	elapsed := now.Sub(l.CreatedAt)

	progress := int64(elapsed+saved) * ProgressScale / int64(originalTotal)
	return clampProgress(progress)
}

// ActualDeliveryAt 返回叠加加速与减时之后的实际预计送达时间。
func (l *Letter) ActualDeliveryAt() time.Time {
	originalTotal := l.ExpectedDeliveryAt.Sub(l.CreatedAt)
	if originalTotal <= 0 {
		return l.ExpectedDeliveryAt
	}
	return l.CreatedAt.Add(l.acceleratedTotal(originalTotal))
}

// acceleratedTotal 按减时道具与加速倍率折算后的总时长。
func (l *Letter) acceleratedTotal(originalTotal time.Duration) time.Duration {
	rate := l.SpeedRate
	if rate <= 0 {
		rate = 1
	}

	total := time.Duration(float64(originalTotal-l.ReduceTime) / rate)
	if total < 0 {
		total = 0
	}
	return total
}

func clampProgress(p int64) int64 {
	switch {
	case p < 0:
		return 0
	case p > ProgressScale:
		return ProgressScale
	default:
		return p
	}
}
