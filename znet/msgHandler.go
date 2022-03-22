package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

//消息处理模块的实现
type MsgHandler struct {
	//存放每个msgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的Worker数量
	WorkerPoolSize uint32
}

// NewMsgHandler 初始化
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取，或在配置文件中由用户设置
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//1 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " not found! Need register!")
		return
	}
	//2 根据msgID调度对应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前msgID是否已存在绑定API处理方法
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!")
}

// StartWorkerPool 启动一个Worker工作池(开启工作池的动作只能发生一次，一个Zinx框架只能有一个工作池，对外暴露)
func (mh *MsgHandler) StartWorkerPool() {
	//根据WorkerPoolSize来分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个Worker被启动
		//1 给当前的Worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is starting...")
	//不断阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的request，执行当前request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 将消息平均分配给不同的Worker
	//根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(),
		" to WorkerID = ", workerID)
	//2 将消息发送给对应的Worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
