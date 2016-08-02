/**
 * Created by banxia on 16/5/17.
 */


var JobInfoList = {

    init:function(){
        JobInfoList.ctrl.initEvent();
    },
    ctrl:{

        initEvent:function(){

            $('.btn_delete').on('click',function(){


                var va = $(this).attr("att");

                layer.confirm('您真的想删除此job吗?', {icon: 3, title:'消息提示'}, function(index){
                    layer.close(index);
                    JobInfoList.ctrl.deleteJobInfo(va);
                });


            });

            $("input[type='checkbox']").on('switchChange.bootstrapSwitch', function (event, state) {
                var mySwitch =  $(this);
                var jobId = $(this).val();
                mySwitch.bootstrapSwitch('state', !state, true);
               if(state) {
                   layer.confirm('您真的想激活此job吗?', {icon: 3, title:'消息提示'}, function(index){
                       layer.close(index);
                       JobInfoList.ctrl.activeJob(mySwitch,1,jobId);
                       mySwitch.bootstrapSwitch('state', state, true);
                   });
               } else {

                   layer.confirm('您真的想注销此job吗?', {icon: 3, title:'消息提示'}, function(index){
                       layer.close(index);
                       JobInfoList.ctrl.activeJob(mySwitch,0,jobId);
                       mySwitch.bootstrapSwitch('state', state, true);
                   });
               }

            });
            $("input[type='checkbox']").bootstrapSwitch();



        },

        // active =1 active 0
        activeJob:function(ele,active,jobId){

            var dat = {"Id":jobId,"active":active};
            layer.load(1);
            $.ajax({
                url:'/jobinfo/active',
                data:dat,
                dataType:'json',
                error:function(XMLHttpRequest, textStatus, errorThrown){
                    layer.closeAll();
                    layer.msg('提交失败', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function(){
                    });

                },
                success:function(data, textStatus, jqXHR){
                    layer.closeAll();
                    if(data.success == true) {
                        if (data.data == 1) {
                            ele.bootstrapSwitch('state', true, true);
                        } else{
                            ele.bootstrapSwitch('state', false, true);
                        }
                        layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                        });
                    }

                },
                type:'POST',
                cache:true

            });
        },
        deleteJobInfo:function(id){

            var dat = {"Id":id};
            layer.load(1);
            $.ajax({
                url:'/jobinfo/delete',
                data:dat,
                dataType:'json',
                error:function(XMLHttpRequest, textStatus, errorThrown){
                    layer.closeAll();
                    layer.msg('提交失败', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function(){
                    });

                },
                success:function(data, textStatus, jqXHR){
                    layer.closeAll();
                    if(data.success == true) {
                        layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                            window.location.href ="/jobinfo/list";
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                        });
                    }

                },
                type:'POST',
                cache:true

            });
        }
    }
};


$(function(){

    JobInfoList.init();
});
