import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;

import javax.xml.soap.Text;

import org.antlr.v4.runtime.tree.TerminalNode;

import parser.CspBaseVisitor;
import parser.CspParser;
import parser.CspParser.AnyContext;
import parser.CspParser.CatchContext;
import parser.CspParser.DEFINITIONContext;
import parser.CspParser.ECHOICEContext;
import parser.CspParser.EventContext;
import parser.CspParser.IDContext;
import parser.CspParser.INTLContext;
import parser.CspParser.ProcContext;
import parser.CspParser.SendContext;

public class CspVisitorInteger extends CspBaseVisitor<Integer>{
	public Map<Integer, String> procs;
	public Map<Integer, String> msgChannels;
	public Map<String, Integer> reverseProcs;
	public Map<Integer, Integer> baseProcsMap;
	public Map<Integer, Integer> rootProcsMap;
	
	public Map<Integer, Integer> nextMap;
	public Map<Integer, Integer> endMap;
	public Map<Integer, ArrayList<Integer>> initMap;
	public Map<Integer, ArrayList<Integer>> msgProcSendMap;
	public Map<Integer, ArrayList<Integer>> msgProcCatchMap;
	public Integer procCntInteger;

	
	public CspVisitorInteger() {
		this(0);
	}
	
	public CspVisitorInteger(Integer startCounter) {
		reverseProcs = new HashMap<String, Integer>();
		procs = new HashMap<Integer, String>();
		msgChannels = new HashMap<Integer, String>();
		baseProcsMap = new HashMap<Integer, Integer>();
		rootProcsMap = new HashMap<Integer, Integer>();
		
		nextMap = new HashMap<Integer, Integer>();
		endMap = new HashMap<Integer, Integer>();
		initMap = new HashMap<Integer, ArrayList<Integer>>();
		msgProcSendMap = new HashMap<Integer, ArrayList<Integer>>();;
		msgProcCatchMap = new HashMap<Integer, ArrayList<Integer>>();;
		
		procCntInteger = 0;
	}
	
    /**
     * 对传入的proc进行处理，将复合的proc分解成独立的proc组合
     * @param ctx 
     */
	private int manageProc(String ctx) {
		return manageProc(ctx, false, false);
	}
	private int manageMsgChannel(String ctx) {
		return manageProc(ctx, false, true);
	}
    private int manageProc(String ctx, boolean isEvent, boolean isMsg) {
    	// 首先识别a -> B 这种原子形式
//    	if (ctx.getChildCount() == 3 && ctx.children.get(0).getChildCount() == 0 
//    				&& ctx.children.get(1).getText() == CspParser.ruleNames[CspParser.ARROW]
//    						&& (ctx.children.get(2).getChild(0).getText() == CspParser.ruleNames[CspParser.Skip]
//    								|| ctx.children.get(2).getChild(0).getText() == CspParser.ruleNames[CspParser.Stop])) {
//			procs.put(procCntInteger, ctx.getText());
//			return procCntInteger++;
//    	}
//    	return -1;
    	if (reverseProcs.containsKey(ctx)) {
    		Integer idInteger = reverseProcs.get(ctx);
    		if (isEvent) {
    			baseProcsMap.put(idInteger, 1);
    		}
    		return reverseProcs.get(ctx);
    	}
    	if (isMsg) {
    		msgChannels.put(procCntInteger, ctx);
    	} else {
    		procs.put(procCntInteger, ctx);
    	}
    	reverseProcs.put(ctx, procCntInteger);
    	if (isEvent) {
    		baseProcsMap.put(procCntInteger, 1);
    	}
		return procCntInteger++;
    }
    
    private void manageMsgChannelProc(Integer chan, Integer proc, boolean isSend) {
		if (isSend) {
			if (!msgProcSendMap.containsKey(chan)) {
	    		ArrayList<Integer> tmpArrayList = new ArrayList<Integer>();
	    		msgProcSendMap.put(chan, tmpArrayList);
			}
			msgProcSendMap.get(chan).add(proc);
		} else {
			if (!msgProcCatchMap.containsKey(chan)) {
	    		ArrayList<Integer> tmpArrayList = new ArrayList<Integer>();
	    		msgProcCatchMap.put(chan, tmpArrayList);
			}
			msgProcCatchMap.get(chan).add(proc);
		}
	}
    
    private void manageNext(Integer first, Integer second) {
    	nextMap.put(first, second);
    }

    private void manageEnd(Integer first, Integer second) {
    	endMap.put(first, second);
    }

    private void manageInit(Integer first, Integer second) {
    	if (!initMap.containsKey(first)) {
    		ArrayList<Integer> tmpArrayList = new ArrayList<Integer>();
    		initMap.put(first, tmpArrayList);
    	}
    	initMap.get(first).add(second);
    }
    
    
    @Override
    public Integer visitDEFINITION(DEFINITIONContext ctx) {
    	visitChildren(ctx);
    	AnyContext def = ctx.any();
    	String textString = def.getText();
    	textString = textString.substring(0, textString.length()-1);
    	Integer defIDInteger = manageProc(textString);
    	rootProcsMap.put(defIDInteger, 1);
    	
    	return defIDInteger;
    }
    
    @Override
    public Integer visitSEMICOL(CspParser.SEMICOLContext ctx) {
    	visitChildren(ctx);
    	ProcContext LeftProc = ctx.proc(0);
    	ProcContext rightProc = ctx.proc(1);
    	
    	Integer procID = manageProc(ctx.getText());
    	Integer leftID = manageProc(LeftProc.getText());
    	Integer rightID = manageProc(rightProc.getText());
    	
    	manageNext(leftID, rightID);
    	manageEnd(rightID, procID);
    	manageInit(procID, leftID);
		return procID;
    }
    
    @Override
    public Integer visitARROW(CspParser.ARROWContext ctx) {
    	visitChildren(ctx);

    	ProcContext rightProc = ctx.proc();
    	EventContext left = ctx.event();
    	Integer procID = manageProc(ctx.getText());
    	// 左边的event是msg ?/! Proc, 这两个都是ID，我们只取右边的
    	Integer leftID = manageProc(left.getChild(2).getText(), true, false);
    	if (!(rightProc.getText().equals("SKIP")
    			|| rightProc.getText().equals("STOP"))) {
	    	Integer rightProcID = manageProc(rightProc.getText());
	    	
	    	manageNext(leftID, rightProcID);
	    	manageEnd(rightProcID, procID);
	    	manageInit(procID, leftID);
    	} else {
    		manageInit(procID, leftID);
    		manageEnd(leftID, procID);
    	}
    	
    	
    	return procID;
    }
    
    @Override
    public Integer visitINTL(INTLContext ctx) {
    	visitChildren(ctx);
    	ProcContext LeftProc = ctx.proc(0);
    	ProcContext rightProc = ctx.proc(1);
    	
    	Integer procID = manageProc(ctx.getText());
    	Integer leftID = manageProc(LeftProc.getText());
    	Integer rightID = manageProc(rightProc.getText());
    	
    	manageEnd(leftID, procID);
    	manageEnd(rightID, procID);
    	manageInit(procID, leftID);
    	manageInit(procID, rightID);
    	
    	return procID;
    }
    
    @Override
    public Integer visitECHOICE(ECHOICEContext ctx) {
    	visitChildren(ctx);
    	
    	ProcContext LeftProc = ctx.proc(0);
    	ProcContext rightProc = ctx.proc(1);
    	
    	Integer procID = manageProc(ctx.getText());
    	Integer leftID = manageProc(LeftProc.getText());
    	Integer rightID = manageProc(rightProc.getText());
    	manageEnd(leftID, procID);
    	manageEnd(rightID, procID);
    	manageInit(procID, leftID);
    	manageInit(procID, rightID);
    	// TODO: 任意一个开始执行后，其他的都不可再执行
    	
    	return procID;
    }
	
	@Override
	public Integer visitParen(CspParser.ParenContext ctx) {
		visitChildren(ctx);
		
    	Integer procID = manageProc(ctx.getText());
    	Integer subProcID = manageProc(ctx.proc().getText());

    	manageInit(procID, subProcID);
    	manageEnd(subProcID, procID);
    	return procID;
	}
	
	@Override
	public Integer visitCatch(CatchContext ctx) {
		visitChildren(ctx);
		// 记录消息通道和proc的关系
    	Integer msgChannelID = manageMsgChannel(ctx.ID(0).getText());
    	Integer procID = manageProc(ctx.ID(1).getText(), true, false);
    	manageMsgChannelProc(msgChannelID, procID, false);

    	Integer totalID = manageProc(ctx.getText());
    	manageInit(totalID, procID);
    	manageEnd(procID, totalID);
    	
    	return -1;
	}
	
	@Override
	public Integer visitSend(SendContext ctx) {
		visitChildren(ctx);
		// 记录消息通道和proc的关系
    	Integer msgChannelID = manageMsgChannel(ctx.ID(0).getText());
    	Integer procID = manageProc(ctx.ID(1).getText(), true, false);
    	manageMsgChannelProc(msgChannelID, procID, true);
    	
    	Integer totalID = manageProc(ctx.getText());
    	manageInit(totalID, procID);
    	manageEnd(procID, totalID);
    	
    	return -1;
	}
}
