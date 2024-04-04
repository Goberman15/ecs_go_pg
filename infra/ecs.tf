resource "aws_ecs_cluster" "ecs_cluster" {
  name = "tf_ecs_go_pg"
}

resource "aws_ecs_task_definition" "web_backend_task" {
  cpu                      = 1024
  memory                   = 4096
  family                   = "web_backend_task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = var.ecs_exec_role_arn

  container_definitions = jsonencode([
    {
      name      = "web_backend"
      cpu       = 1024
      memory    = 4096
      image     = var.docker_image_url
      essential = true
      portMappings = [
        {
          containerPort = 8088
          hostPort      = 8088
        }
      ]
      environment = [
        { "name" : "DB_HOST", "value" : "${var.db_host}" },
        { "name" : "DB_USER", "value" : "${var.db_user}" },
        { "name" : "DB_NAME", "value" : "${var.db_name}" },
        { "name" : "DB_PASS", "value" : "${var.db_pass}" },
        { "name" : "DB_PORT", "value" : "${var.db_port}" },
        { "name" : "DB_PARAMS", "value" : "${var.db_params}" }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-create-group  = "true"
          awslogs-group         = "/ecs/web_backend_task"
          awslogs-region        = "ap-southeast-1"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_service" "web_backend_service" {
  name            = "web_backend_service"
  cluster         = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.web_backend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.subnet_ids
    security_groups  = [var.sg_id]
    assign_public_ip = true
  }
}
