/**
 * Created by banxia on 16/5/16.
 */


var AddJob = {

        init:function(){

            AddJob.ctrl.initEvent();
            AddJob.ctrl.initValidate();
        },
    ctrl:{
        initEvent:function(){

        },
        initValidate:function(){

            $('#addForm').bootstrapValidator({
                message: '这个值不可用',
                feedbackIcons: {
                    valid: 'glyphicon glyphicon-ok',
                    invalid: 'glyphicon glyphicon-remove',
                    validating: 'glyphicon glyphicon-refresh'
                },
                fields:{

                    Name:{
                        message: 'Job名称不可用',
                        validators: {
                            notEmpty: {
                                message: 'job名称不能为空'
                            },
                            stringLength: {
                                min: 3,
                                max: 30,
                                message: 'job应在3~30个字符之间'
                            }

                        }
                    },
                    Group:{
                        message: '任务分组不可用',
                        validators: {
                            notEmpty: {
                                message: '任务分组不能为空'
                            },
                            stringLength: {
                                min: 3,
                                max: 30,
                                message: '任务分组应在3~30个字符之间'
                            }

                        }
                    },
                    Urls:{
                        message: '目标服务器地址不可用',
                        validators: {
                            notEmpty: {
                                message: '目标服务器地址不能为空'
                            }

                        }
                    },
                    Cron:{
                        message: 'Cron表达式不可用',
                        validators: {
                            notEmpty: {
                                message: 'Cron表达式不能为空'
                            }

                        }
                    }

                }
            }).on('success.form.bv', function(e) {
                // Prevent form submission
                e.preventDefault();

                // Get the form instance
                var $form = $(e.target);
                // Get the BootstrapValidator instance
                var bv = $form.data('bootstrapValidator');

                var dat = $form.serialize();
                layer.load(1);
                $.ajax({
                    url:'/jobinfo/add',
                    data:dat,
                    dataType:'json',
                    error:function(XMLHttpRequest, textStatus, errorThrown){
                        layer.closeAll();
                        $form
                            .bootstrapValidator('disableSubmitButtons', false)
                    },
                    success:function(data, textStatus, jqXHR){
                        layer.closeAll();
                        $form
                            .bootstrapValidator('disableSubmitButtons', false)
                        if(data.success == true) {
                            layer.msg(data.message, {
                                icon: 1,
                                time: 2000 //2秒关闭（如果不配置，默认是3秒）
                            }, function(){
                                window.location.href ="/jobinfo/list";
                            });

                        } else {
                            layer.msg(data.Message, {
                                icon: 5,
                                time: 2000 //2秒关闭（如果不配置，默认是3秒）
                            }, function(){
                            });
                        }

                    },
                    type:'POST',
                    cache:true

                });


            });
        }
    }
};
$(function(){


     AddJob.init();

});
