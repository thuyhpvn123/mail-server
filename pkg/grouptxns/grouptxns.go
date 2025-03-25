package grouptxns

import (
	"fmt"
	"time"

	"gomail/types"

	"github.com/ethereum/go-ethereum/common"
)

// Item đại diện cho một phần tử cần được nhóm (transaction)
type Item struct {
	ID        int
	Array     []common.Address
	GroupID   int
	Tx        types.Transaction
	TimeStart time.Time
}

// UnionFind là cấu trúc dữ liệu Union-Find
type UnionFind struct {
	parent []int
	rank   []int
}

// GroupResult đại diện cho kết quả xử lý một nhóm giao dịch.
type GroupResult struct {
	Transactions     []types.Transaction
	Receipts         []types.Receipt
	ExecuteSCResults []types.ExecuteSCResult
	Error            error
	AsRoot           common.Hash
}

// RelativeGroup đại diện cho một nhóm giao dịch liên quan
type RelativeGroup struct {
	GroupID   int
	Items     []Item
	Relatives []common.Address
}

// TotalGas tính tổng gas của tất cả các item trong nhóm
func (rg *RelativeGroup) TotalGas() uint64 {
	totalGas := uint64(0)
	for _, item := range rg.Items {
		totalGas += item.Tx.MaxGas()
	}
	return totalGas
}

// TotalTime tính tổng thời gian của tất cả các item trong nhóm
func (rg *RelativeGroup) TotalTime() uint64 {
	totalTime := uint64(0)
	for _, item := range rg.Items {
		totalTime += item.Tx.MaxTimeUse()
	}
	return totalTime
}

// NewUnionFind tạo một đối tượng UnionFind mới
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent, rank}
}

// Find tìm cha của một phần tử
func (uf *UnionFind) Find(i int) int {
	if uf.parent[i] == i {
		return i
	}
	uf.parent[i] = uf.Find(uf.parent[i])
	return uf.parent[i]
}

// Union hợp nhất hai phần tử
func (uf *UnionFind) Union(i, j int) {
	rootI := uf.Find(i)
	rootJ := uf.Find(j)
	if rootI != rootJ {
		if uf.rank[rootI] < uf.rank[rootJ] {
			uf.parent[rootI] = rootJ
		} else if uf.rank[rootI] > uf.rank[rootJ] {
			uf.parent[rootJ] = rootI
		} else {
			uf.parent[rootJ] = rootI
			uf.rank[rootI]++
		}
	}
}
func GroupItems(items []Item) ([]Item, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("input slice is empty")
	}

	// Sử dụng map để ánh xạ địa chỉ -> danh sách các chỉ số của các item chứa địa chỉ đó
	addressToIndices := make(map[string][]int)
	for i, item := range items {
		for _, addr := range item.Array {
			addressToIndices[addr.Hex()] = append(addressToIndices[addr.Hex()], i)
		}
	}

	// Khởi tạo Union-Find
	uf := NewUnionFind(len(items))

	// Liên kết các phần tử có địa chỉ chung
	for _, indices := range addressToIndices {
		for i := 1; i < len(indices); i++ {
			uf.Union(indices[0], indices[i])
		}
	}

	// Gán GroupID cho mỗi item dựa trên Union-Find
	groupIDMap := make(map[int]int)
	groupIDCounter := 1
	for i := range items {
		root := uf.Find(i)
		if _, ok := groupIDMap[root]; !ok {
			groupIDMap[root] = groupIDCounter
			groupIDCounter++
		}
		items[i].GroupID = groupIDMap[root]
	}

	return items, nil
}

// groupByGroupID nhóm các phần tử theo GroupID
func GroupByGroupID(items []Item) [][]Item {
	groups := make(map[int][]Item)
	for _, item := range items {
		groups[item.GroupID] = append(groups[item.GroupID], item)
	}

	result := make([][]Item, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	return result
}

func GroupAndLimitTransactionsOptimized(items []Item, maxGroupGas uint64, maxTotalGas uint64, maxGroupTimes uint64, maxTotalTime uint64) ([]RelativeGroup, []Item, error) {

	relativeGroups := []RelativeGroup{}
	excludedItems := []Item{}
	totalGas := uint64(0)
	totalTime := uint64(0)

	// Map ánh xạ địa chỉ tới groupID
	addressToGroup := make(map[string]int)

	for _, item := range items {
		// Nếu giao dịch chỉ đọc, tạo nhóm riêng
		if item.Tx.GetReadOnly() {
			newGroup := RelativeGroup{
				GroupID: len(relativeGroups),
				Items:   []Item{item},
			}
			relativeGroups = append(relativeGroups, newGroup)
			continue // Tiếp tục vòng lặp cho item tiếp theo
		}

		var selectedGroup *RelativeGroup

		// Kiểm tra các nhóm có thể thêm item
		for _, addr := range item.Array {
			if groupID, exists := addressToGroup[addr.Hex()]; exists {
				group := &relativeGroups[groupID]
				newGas := group.TotalGas() + item.Tx.MaxGas()
				newTime := group.TotalTime() + item.Tx.MaxTimeUse()

				// Nếu nhóm này đủ điều kiện thì chọn
				if newGas <= maxGroupGas && newTime <= maxGroupTimes {
					selectedGroup = group
					break
				}
			}
		}

		if selectedGroup != nil {
			// Nếu tìm thấy nhóm phù hợp, thêm item vào nhóm
			selectedGroup.Items = append(selectedGroup.Items, item)
			for _, addr := range item.Array {
				addressToGroup[addr.Hex()] = selectedGroup.GroupID
			}
		} else {
			// Nếu không tìm thấy nhóm phù hợp, tạo nhóm mới
			newGroup := RelativeGroup{
				GroupID: len(relativeGroups),
				Items:   []Item{item},
			}
			newGas := item.Tx.MaxGas()
			newTime := item.Tx.MaxTimeUse()

			// Kiểm tra các điều kiện giới hạn
			if newGas <= maxGroupGas && newTime <= maxGroupTimes &&
				totalGas+newGas <= maxTotalGas && totalTime+newTime <= maxTotalTime && len(relativeGroups) < 5000 {
				relativeGroups = append(relativeGroups, newGroup)
				for _, addr := range item.Array {
					addressToGroup[addr.Hex()] = newGroup.GroupID
				}
				totalGas += newGas
				totalTime += newTime
			} else {
				// Thêm vào danh sách loại bỏ nếu không hợp lệ
				excludedItems = append(excludedItems, item)
			}
		}
	}

	return relativeGroups, excludedItems, nil
}
